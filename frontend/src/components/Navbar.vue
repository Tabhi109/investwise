<template>
  <header class="navbar glass-card">
    <div class="nav-container">
      <!-- Brand Logo -->
      <div class="brand" @click="$emit('update:activeTab', 'portfolio')">
        <div class="brand-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="20" x2="12" y2="10"></line>
            <line x1="18" y1="20" x2="18" y2="4"></line>
            <line x1="6" y1="20" x2="6" y2="16"></line>
          </svg>
        </div>
        <div class="brand-text">
          <span class="brand-name brand-font">InvestWise</span>
          <span class="brand-tagline">Smart Investing Platform</span>
        </div>
      </div>

      <!-- Desktop Nav Tabs -->
      <div class="desktop-tabs">
        <button 
          class="desktop-tab" 
          :class="{ 'active': activeTab === 'portfolio' }"
          @click="$emit('update:activeTab', 'portfolio')"
        >
          Portfolio
        </button>
        <button 
          class="desktop-tab" 
          :class="{ 'active': activeTab === 'markets' }"
          @click="$emit('update:activeTab', 'markets')"
        >
          Markets & Trade
        </button>
        <button 
          class="desktop-tab" 
          :class="{ 'active': activeTab === 'risk' }"
          @click="$emit('update:activeTab', 'risk')"
        >
          Quant Coach
        </button>
        <button 
          class="desktop-tab" 
          :class="{ 'active': activeTab === 'learn' }"
          @click="$emit('update:activeTab', 'learn')"
        >
          <span class="tab-spark">✨</span> Academy
        </button>
      </div>

      <!-- Right Actions -->
      <div class="nav-actions">
        <span class="badge badge-success ws-status">
          <span class="pulse-dot"></span> <span class="ws-text">Live Feed</span>
        </span>

        <div v-if="user" class="user-info">
          <div class="user-avatar" :title="user.email">{{ user.email.charAt(0).toUpperCase() }}</div>
          <span class="user-email">{{ user.email }}</span>
          <button @click="$emit('logout')" class="btn btn-secondary btn-sm" id="logout-btn" title="Sign out">
            Logout
          </button>
        </div>
        <div v-else>
          <button @click="$emit('open-auth')" class="btn btn-primary btn-sm" id="login-modal-btn">
            Sign In / Demo
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
defineProps({
  user: Object,
  activeTab: {
    type: String,
    default: 'portfolio'
  }
})
defineEmits(['open-auth', 'logout', 'update:activeTab'])
</script>

<style scoped>
.navbar {
  margin-bottom: 1.25rem;
  padding: 0.75rem 1.25rem;
  border-radius: 0 0 16px 16px;
  position: sticky;
  top: 0;
  z-index: 800;
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
  gap: 0.65rem;
  cursor: pointer;
}

.brand-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--primary), var(--accent));
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-name {
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1.1;
}

.brand-tagline {
  font-size: 0.65rem;
  color: var(--text-dim);
  font-weight: 500;
}

@media (max-width: 480px) {
  .brand-tagline {
    display: none;
  }
}

/* Desktop Tabs */
.desktop-tabs {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(0, 0, 0, 0.3);
  padding: 0.25rem 0.35rem;
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

@media (max-width: 768px) {
  .desktop-tabs {
    display: none;
  }
}

.desktop-tab {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 0.825rem;
  font-weight: 600;
  padding: 0.45rem 0.9rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
}

.desktop-tab:hover {
  color: white;
}

.desktop-tab.active {
  background: rgba(99, 102, 241, 0.2);
  color: #c7d2fe;
  border: 1px solid rgba(99, 102, 241, 0.4);
}

.tab-spark {
  font-size: 0.85rem;
}

/* Right Actions */
.nav-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.ws-status {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.7rem;
}

@media (max-width: 480px) {
  .ws-text {
    display: none;
  }
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
  gap: 0.6rem;
}

.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--accent));
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.8rem;
}

.user-email {
  font-size: 0.825rem;
  color: var(--text-muted);
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 640px) {
  .user-email {
    display: none;
  }
}

.btn-sm {
  padding: 0.4rem 0.85rem;
  font-size: 0.8rem;
  min-height: 36px;
}
</style>
