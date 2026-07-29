<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content glass-card">
      <div class="modal-header">
        <h3 class="modal-title">{{ isRegister ? 'Create Account' : 'Welcome Back' }}</h3>
        <button class="close-btn" @click="$emit('close')">&times;</button>
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

        <button type="submit" class="btn btn-primary btn-block" :disabled="loading" id="auth-submit-btn">
          {{ loading ? 'Processing...' : (isRegister ? 'Register & Claim $100k Cash' : 'Sign In') }}
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
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 100%;
  max-width: 420px;
  padding: 2rem;
  position: relative;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
}

.modal-title {
  font-size: 1.4rem;
  font-weight: 700;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

.tab-buttons {
  display: flex;
  background: rgba(0, 0, 0, 0.4);
  padding: 4px;
  border-radius: 10px;
  margin-bottom: 1.5rem;
}

.tab-btn {
  flex: 1;
  padding: 0.5rem;
  border: none;
  background: none;
  color: var(--text-muted);
  font-weight: 600;
  font-size: 0.875rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: var(--primary);
  color: white;
}

.form-group {
  margin-bottom: 1.15rem;
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
  padding: 0.75rem;
  margin-top: 0.5rem;
}
</style>
