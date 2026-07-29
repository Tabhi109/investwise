<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content glass-card">
      <div class="modal-header">
        <h3 class="modal-title">Execute Trade</h3>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>

      <div class="trade-type-toggle">
        <button 
          :class="['toggle-btn', { 'active-buy': tradeType === 'BUY' }]" 
          @click="tradeType = 'BUY'"
        >
          BUY
        </button>
        <button 
          :class="['toggle-btn', { 'active-sell': tradeType === 'SELL' }]" 
          @click="tradeType = 'SELL'"
        >
          SELL
        </button>
      </div>

      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>Select Asset</label>
          <select v-model="selectedTicker" class="form-control" id="trade-ticker-select">
            <option v-for="t in tickers" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>

        <div class="price-box">
          <div class="price-label">Current Market Price</div>
          <div class="price-val">${{ currentPrice.toFixed(2) }}</div>
        </div>

        <div class="form-group">
          <label>Number of Shares</label>
          <input 
            v-model.number="shares" 
            type="number" 
            step="0.0001"
            min="0.0001"
            class="form-control" 
            placeholder="10" 
            required 
            id="trade-shares-input"
          />
        </div>

        <div class="summary-box">
          <div class="summary-row">
            <span>Estimated Total {{ tradeType === 'BUY' ? 'Cost' : 'Proceeds' }}</span>
            <span class="font-bold">${{ estimatedTotal.toFixed(2) }}</span>
          </div>
          <div v-if="tradeType === 'BUY'" class="summary-row subtext">
            <span>Available Cash</span>
            <span>${{ cashBalance.toFixed(2) }}</span>
          </div>
        </div>

        <div v-if="error" class="error-banner">
          {{ error }}
        </div>

        <button 
          type="submit" 
          :class="['btn', 'btn-block', tradeType === 'BUY' ? 'btn-success' : 'btn-danger']"
          :disabled="loading || (tradeType === 'BUY' && estimatedTotal > cashBalance)"
          id="confirm-trade-btn"
        >
          {{ loading ? 'Executing...' : `Confirm ${tradeType} Order` }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  tickers: { type: Array, default: () => ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'NVDA'] },
  priceMap: { type: Object, default: () => ({}) },
  initialTicker: { type: String, default: 'AAPL' },
  initialType: { type: String, default: 'BUY' },
  cashBalance: { type: Number, default: 0 },
  loading: Boolean,
  error: String
})

const emit = defineEmits(['close', 'execute-trade'])

const tradeType = ref(props.initialType)
const selectedTicker = ref(props.initialTicker)
const shares = ref(1)

watch(() => props.initialTicker, (val) => { if (val) selectedTicker.value = val })
watch(() => props.initialType, (val) => { if (val) tradeType.value = val })

const currentPrice = computed(() => {
  const pData = props.priceMap[selectedTicker.value]
  return pData && pData.price ? pData.price : 150.0
})

const estimatedTotal = computed(() => {
  return (shares.value || 0) * currentPrice.value
})

function handleSubmit() {
  emit('execute-trade', {
    ticker: selectedTicker.value,
    type: tradeType.value,
    shares: shares.value
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
  max-width: 440px;
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
  font-size: 1.35rem;
  font-weight: 700;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

.trade-type-toggle {
  display: flex;
  background: rgba(0, 0, 0, 0.4);
  padding: 4px;
  border-radius: 10px;
  margin-bottom: 1.25rem;
}

.toggle-btn {
  flex: 1;
  padding: 0.6rem;
  border: none;
  background: none;
  color: var(--text-muted);
  font-weight: 700;
  font-size: 0.9rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.toggle-btn.active-buy {
  background: #10b981;
  color: white;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.toggle-btn.active-sell {
  background: #f43f5e;
  color: white;
  box-shadow: 0 4px 12px rgba(244, 63, 94, 0.4);
}

.price-box {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-color);
  padding: 0.85rem;
  border-radius: 10px;
  margin-bottom: 1.15rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price-label {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.price-val {
  font-family: var(--font-heading);
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--primary);
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

.summary-box {
  background: rgba(15, 23, 42, 0.9);
  border: 1px dashed var(--border-color);
  padding: 0.85rem 1rem;
  border-radius: 10px;
  margin-bottom: 1.25rem;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
}

.summary-row.subtext {
  font-size: 0.75rem;
  color: var(--text-dim);
  margin-top: 0.35rem;
}

.font-bold {
  font-weight: 700;
  color: white;
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
}
</style>
