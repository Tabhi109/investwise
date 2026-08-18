# 📈 InvestWise

A full-stack real-time portfolio management and quantitative risk analysis platform.

## 🚀 Quick Start

### Development Mode
```bash
docker compose up -d --build
```
Access the application at [http://localhost:8000](http://localhost:8000).

### Production Deployment
For complete production deployment instructions (Docker Compose on VPS / Cloud VMs, Railway, Render, SSL/TLS, and database backups), see [DEPLOYMENT.md](file:///Users/abhi/Desktop/Projects/investwise/DEPLOYMENT.md).

```bash
cp .env.example .env
# Edit .env with your production credentials
docker compose -f docker-compose.prod.yml up -d --build
```

## 🛠️ Tech Stack
- **Backend**: Go 1.22+, Gorilla WebSocket, pgx, go-redis
- **Frontend**: Vue 3, Vite, Chart.js / Canvas
- **Database**: PostgreSQL 16 (with automated embedded migrations)
- **Cache & Broker**: Redis 7
- **Reverse Proxy**: Nginx