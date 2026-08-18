# 🚀 InvestWise Deployment Guide

This guide covers production deployment workflows for **InvestWise**, ranging from single-server Docker Compose setups (AWS EC2, DigitalOcean, Hetzner, Linode) to managed cloud container platforms (Render, Railway, Fly.io, AWS ECS).

---

## 🏗️ Architecture Overview

InvestWise consists of 5 coordinated services:
- **Go Backend API & Real-time Engine (`app`)**: Handles REST endpoints, WebSocket price feeds, and async risk calculation workers. Runs on port `8080`.
- **Vue 3 Frontend (`frontend`)**: Single Page Application built with Vite and served statically via Nginx on port `80`.
- **Nginx Gateway (`nginx`)**: Reverse proxy routing web traffic, REST API calls, and upgrading `/ws` WebSocket connections.
- **PostgreSQL 16 (`postgres`)**: Relational database with automatic migration execution on startup.
- **Redis 7 (`redis`)**: In-memory cache and pub/sub message broker.

---

## ⚡ Option 1: VPS / Cloud VM (Recommended - AWS EC2, DigitalOcean, Hetzner)

This is the most cost-effective and straightforward deployment option.

### Prerequisites
- A Linux server (Ubuntu 22.04 / 24.04 LTS recommended)
- Docker & Docker Compose installed
- A domain name pointing to your server's public IP (optional, for SSL)

### Step 1: Install Docker & Git on Server

```bash
# Update packages
sudo apt update && sudo apt upgrade -y

# Install Docker & Docker Compose plugin
sudo apt install -y git curl ufw
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

*Log out and log back in for docker group permissions to take effect.*

### Step 2: Clone Repository & Configure Environment

```bash
git clone https://github.com/Tabhi109/InvestWise.git
cd InvestWise

# Create production environment configuration
cp .env.example .env
```

Generate secure secrets and update `.env`:

```bash
# Generate a strong JWT Secret
openssl rand -base64 32
```

Edit `.env` using `nano .env`:
```ini
APP_ENV=production
PORT=8080
HTTP_PORT=80

POSTGRES_USER=investwise_admin
POSTGRES_PASSWORD=YOUR_STRONG_PASSWORD_HERE
POSTGRES_DB=investwise

REDIS_URL=redis://redis:6379/0
JWT_SECRET=YOUR_GENERATED_JWT_SECRET_HERE
JWT_TTL_HOURS=24
WORKER_COUNT=8
RISK_FREE_RATE=0.04
MARKET_UPDATE_INTERVAL_MS=1500
```

### Step 3: Launch Production Stack

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

### Step 4: Verify Deployment

```bash
# Check running containers
docker compose -f docker-compose.prod.yml ps

# View live backend logs
docker compose -f docker-compose.prod.yml logs -f app

# Test health / ping endpoint
curl http://localhost/ping
```

### Step 5: Configure Firewall & Free SSL (Let's Encrypt / Certbot)

```bash
# Allow SSH and Web traffic
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

#### Adding SSL with Certbot & Nginx Reverse Proxy:
If you want automatic HTTPS with Let's Encrypt:
```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

---

## ☁️ Option 2: Managed Cloud Platforms (Railway / Render)

If you prefer managed infrastructure without maintaining Linux servers:

### Deploying on Railway:
1. Fork or push the InvestWise repo to GitHub.
2. Log into [Railway.app](https://railway.app).
3. Create a **New Project** and add:
   - **PostgreSQL Database** (Provisioned in 1 click).
   - **Redis Database** (Provisioned in 1 click).
   - **Backend Service**: Connect GitHub repo, set root directory to `backend`.
     - Add environment variables:
       - `DATABASE_URL`: `${{Postgres.DATABASE_URL}}`
       - `REDIS_URL`: `${{Redis.REDIS_URL}}`
       - `JWT_SECRET`: Random 32+ character string
       - `APP_ENV`: `production`
       - `PORT`: `8080`
   - **Frontend Service**: Connect GitHub repo, set root directory to `frontend`.
     - Dockerfile handles build & static serving.

---

## 🛠️ Operations & Maintenance Runbook

### Updating Application to Latest Version
```bash
cd InvestWise
git pull origin main
docker compose -f docker-compose.prod.yml up -d --build
```

### Database Backups (PostgreSQL)
```bash
# Create an on-demand compressed backup
docker exec -t $(docker ps -qf "name=postgres") pg_dump -U investwise_admin investwise | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Automated daily cron backup
(crontab -l 2>/dev/null; echo "0 2 * * * docker exec -t \$(docker ps -qf 'name=postgres') pg_dump -U investwise_admin investwise | gzip > ~/investwise_backup_\$(date +\%Y\%m\%d).sql.gz") | crontab -
```

### Database Restore
```bash
gunzip -c backup_YYYYMMDD_HHMMSS.sql.gz | docker exec -i $(docker ps -qf "name=postgres") psql -U investwise_admin -d investwise
```

### Restarting / Troubleshooting
```bash
# Restart all services
docker compose -f docker-compose.prod.yml restart

# Restart specific service
docker compose -f docker-compose.prod.yml restart app

# View logs for a specific service
docker compose -f docker-compose.prod.yml logs -f app
```

---

## 🔒 Production Hardening Checklist

- [ ] Changed default `POSTGRES_PASSWORD` and `JWT_SECRET` in `.env`.
- [ ] Database ports (`5432` / `6379`) are isolated in private Docker networks and not exposed to the public internet.
- [ ] Enabled HTTPS / SSL certificate.
- [ ] Configured automatic database backups.
- [ ] Verified WebSocket connectivity over secure protocols (`wss://`).
