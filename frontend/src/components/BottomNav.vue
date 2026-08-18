<template>
  <nav class="bottom-nav glass-card" aria-label="Mobile Navigation">
    <button 
      class="nav-tab" 
      :class="{ 'active': activeTab === 'portfolio' }"
      @click="$emit('update:activeTab', 'portfolio')"
      id="tab-portfolio"
    >
      <div class="tab-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 3v18h18"/>
          <path d="m19 9-5 5-4-4-3 3"/>
        </svg>
      </div>
      <span class="tab-label">Portfolio</span>
    </button>

    <button 
      class="nav-tab" 
      :class="{ 'active': activeTab === 'markets' }"
      @click="$emit('update:activeTab', 'markets')"
      id="tab-markets"
    >
      <div class="tab-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 14 14"/>
        </svg>
      </div>
      <span class="tab-label">Markets</span>
    </button>

    <button 
      class="nav-tab" 
      :class="{ 'active': activeTab === 'risk' }"
      @click="$emit('update:activeTab', 'risk')"
      id="tab-risk"
    >
      <div class="tab-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
        </svg>
      </div>
      <span class="tab-label">Quant Coach</span>
    </button>

    <button 
      class="nav-tab" 
      :class="{ 'active': activeTab === 'learn' }"
      @click="$emit('update:activeTab', 'learn')"
      id="tab-learn"
    >
      <div class="tab-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
          <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>
        </svg>
      </div>
      <span class="tab-label">Academy</span>
    </button>
  </nav>
</template>

<script setup>
defineProps({
  activeTab: {
    type: String,
    default: 'portfolio'
  }
})

defineEmits(['update:activeTab'])
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0.75rem;
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 1.5rem);
  max-width: 480px;
  height: var(--bottom-nav-height);
  background: rgba(15, 23, 42, 0.88);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 20px;
  display: flex;
  justify-content: space-around;
  align-items: center;
  z-index: 900;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.6);
  padding: 0 0.5rem;
}

@media (min-width: 769px) {
  .bottom-nav {
    display: none;
  }
}

.nav-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: var(--text-dim);
  padding: 0.4rem 0.2rem;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  gap: 0.2rem;
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}

.tab-label {
  font-size: 0.68rem;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.nav-tab.active {
  color: #a5b4fc;
}

.nav-tab.active .tab-icon {
  transform: translateY(-2px);
  color: #818cf8;
}

.nav-tab.active .tab-label {
  color: white;
  font-weight: 700;
}
</style>
