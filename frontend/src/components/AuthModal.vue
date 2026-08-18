<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content glass-card modal-animate-sheet">
      <div class="modal-header">
        <div class="header-left">
          <span class="chip-sm">ACCOUNT ACCESS</span>
          <h3 class="modal-title">{{ isRegister ? 'Create Demo Account' : 'Welcome Back' }}</h3>
        </div>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>

      <!-- Quick Demo Login Action -->
      <div class="demo-access-card">
        <div class="demo-info">
          <span class="demo-title">⚡ Instant Demo Mode</span>
          <span class="demo-sub">Zero setup required — get $100,000 in simulated funds instantly.</span>
        </div>
        <button type="button" class="btn btn-accent btn-sm" @click="handleQuickDemo">
          Try Live Demo
        </button>
      </div>

      <div class="divider">
        <span>or sign in with email</span>
      </div>

      <div class="tab-buttons">
        <button 
          :class="['tab-btn', { active: !isRegister }]" 
          @click="isRegister = false"
        >
          Sign In
        </button>
        <button 
          :class="['tab-btn', { active: isRegister }]" 
          @click="isRegister = true"
        >
          Register
        </button>
      </div>

      <form @submit.prevent="handleSubmit">
        <div v-if="isRegister" class="form-group">
          <label>Username</label>
          <input 
            v-model="form.username" 
            type="text" 
            class="form-control" 
            placeholder="trader123" 
            required 
            minlength="3"
            id="auth-username"
          />
        </div>

        <div class="form-group">
          <label>Email Address</label>
          <input 
            v-model="form.email" 
            type="email" 
            class="form-control" 
            placeholder="trader@investwise.io" 
            required 
            id="auth-email"
          />
        </div>

        <div class="form-group">
          <label>Password</label>
          <input 
            v-model="form.password" 
            type="password" 
            class="form-control" 
            placeholder="••••••••" 
            required 
            minlength="6"
            id="auth-password"
          />
        </div>

        <div v-if="error" class="error-banner">
          {{ error }}
        </div>

        <button type="submit" class="btn btn-primary btn-block btn-lg" :disabled="loading" id="auth-submit-btn">
          {{ loading ? 'Processing...' : (isRegister ? 'Register & Claim $100k Cash' : 'Sign In to Account') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const props = defineProps({
  loading: Boolean,
  error: String
})

const emit = defineEmits(['close', 'submit-auth'])

const isRegister = ref(false)
const form = reactive({
  username: '',
  email: '',
  password: ''
})

function handleQuickDemo() {
  const randomId = Math.floor(Math.random() * 90000) + 10000
  emit('submit-auth', {
    isRegister: true,
    username: `investor_${randomId}`,
    email: `demo_${randomId}@investwise.io`,
    password: 'password123'
  })
}

function handleSubmit() {
  emit('submit-auth', {
    isRegister: isRegister.value,
    username: form.username,
    email: form.email,
    password: form.password
  })
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  width: 100%;
  max-width: 440px;
  padding: 1.75rem;
  position: relative;
  background: rgba(15, 23, 42, 0.95);
  max-height: 90vh;
  overflow-y: auto;
}

@media (max-width: 640px) {
  .modal-backdrop {
    align-items: flex-end;
    padding: 0;
  }

  .modal-content {
    border-radius: 24px 24px 0 0;
    max-height: 88vh;
    padding: 1.25rem 1.25rem 2rem 1.25rem;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.15rem;
}

.chip-sm {
  font-size: 0.65rem;
  font-weight: 800;
  color: #a5b4fc;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.modal-title {
  font-size: 1.3rem;
  font-weight: 800;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.75rem;
  cursor: pointer;
  line-height: 1;
}

.demo-access-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.15), rgba(6, 182, 212, 0.1));
  border: 1px solid rgba(139, 92, 246, 0.35);
  padding: 0.85rem 1rem;
  border-radius: 12px;
  margin-bottom: 1.25rem;
  gap: 0.75rem;
}

.demo-info {
  display: flex;
  flex-direction: column;
}

.demo-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: #c4b5fd;
}

.demo-sub {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  margin-bottom: 1.25rem;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--border-color);
}

.divider span {
  padding: 0 0.5rem;
  font-size: 0.75rem;
  color: var(--text-dim);
}

.tab-buttons {
  display: flex;
  background: rgba(0, 0, 0, 0.4);
  padding: 4px;
  border-radius: 12px;
  margin-bottom: 1.25rem;
}

.tab-btn {
  flex: 1;
  padding: 0.55rem;
  border: none;
  background: none;
  color: var(--text-muted);
  font-weight: 700;
  font-size: 0.85rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: var(--primary);
  color: white;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  margin-bottom: 0.35rem;
}

.error-banner {
  background: var(--danger-glow);
  border: 1px solid rgba(244, 63, 94, 0.4);
  color: #fb7185;
  padding: 0.6rem 0.85rem;
  border-radius: 8px;
  font-size: 0.825rem;
  margin-bottom: 1.25rem;
}

.btn-block {
  width: 100%;
}

.btn-lg {
  padding: 0.85rem;
  font-size: 0.95rem;
  font-weight: 700;
  border-radius: 12px;
}
</style>
