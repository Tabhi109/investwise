<template>
  <header class="navbar glass-card">
    <div class="nav-container">
      <div class="brand">
        <div class="brand-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="20" x2="12" y2="10"></line>
            <line x1="18" y1="20" x2="18" y2="4"></line>
            <line x1="6" y1="20" x2="6" y2="16"></line>
          </svg>
        </div>
        <span class="brand-name brand-font">InvestWise</span>
        <span class="badge badge-success ws-status">
          <span class="pulse-dot"></span> Live Feed
        </span>
      </div>

      <div class="nav-actions">
        <div v-if="user" class="user-info">
          <div class="user-avatar">{{ user.email.charAt(0).toUpperCase() }}</div>
          <span class="user-email">{{ user.email }}</span>
          <button @click="$emit('logout')" class="btn btn-secondary btn-sm" id="logout-btn">
            Logout
          </button>
        </div>
        <div v-else>
          <button @click="$emit('open-auth')" class="btn btn-primary" id="login-modal-btn">
            Sign In / Register
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
defineProps({
  user: Object
})
defineEmits(['open-auth', 'logout'])
</script>

<style scoped>
.navbar {
  margin-bottom: 1.5rem;
  padding: 0.85rem 1.5rem;
  border-radius: 0 0 16px 16px;
}

.nav-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.brand-icon {
  width: 38px;
  height: 38px;
  background: linear-gradient(135deg, var(--primary), var(--accent));
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.brand-name {
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.ws-status {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.7rem;
}

.pulse-dot {
  width: 6px;
  height: 6px;
  background-color: #10b981;
  border-radius: 50%;
  box-shadow: 0 0 8px #10b981;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% { transform: scale(0.95); opacity: 0.8; }
  50% { transform: scale(1.3); opacity: 1; }
  100% { transform: scale(0.95); opacity: 0.8; }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.85rem;
}

.user-email {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.btn-sm {
  padding: 0.4rem 0.85rem;
  font-size: 0.8rem;
}
</style>
