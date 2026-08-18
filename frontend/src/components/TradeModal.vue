<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content glass-card modal-animate-sheet">
      <div class="modal-header">
        <div class="header-title-wrap">
          <span class="chip-sm">ORDER EXECUTION</span>
          <h3 class="modal-title">Trade Stock</h3>
        </div>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>

      <!-- Trade Type Switcher -->
      <div class="trade-type-toggle">
        <button 
          :class="['toggle-btn', { 'active-buy': tradeType === 'BUY' }]" 
          @click="tradeType = 'BUY'"
          id="toggle-buy-btn"
        >
          BUY ORDER
        </button>
        <button 
          :class="['toggle-btn', { 'active-sell': tradeType === 'SELL' }]" 
          @click="tradeType = 'SELL'"
          id="toggle-sell-btn"
        >
          SELL ORDER
        </button>
      </div>

      <form @submit.prevent="handleSubmit">
        <!-- Asset Selector -->
        <div class="form-group">
          <label>Select Asset</label>
          <div class="select-wrap">
            <select v-model="selectedTicker" class="form-control" id="trade-ticker-select">
              <option v-for="t in tickers" :key="t" :value="t">{{ t }} - {{ getCompanyName(t) }}</option>
            </select>
          </div>
        </div>

        <!-- Live Price & Quote -->
        <div class="price-quote-card">
          <div class="quote-info">
            <span class="quote-label">Market Price</span>
            <span class="quote-price font-heading">${{ currentPrice.toFixed(2) }}</span>
          </div>
          <span class="live-dot-badge"><span class="pulse-dot"></span> Real-time</span>
        </div>

        <!-- Number of Shares & Presets -->
        <div class="form-group">
          <div class="label-row">
            <label>Number of Shares</label>
            <span v-if="tradeType === 'SELL'" class="owned-shares-hint">
              Owned: {{ ownedShares.toFixed(2) }}
            </span>
          </div>
          <input 
            v-model.number="shares" 
            type="number" 
            step="0.0001"
            min="0.0001"
            class="form-control shares-input" 
            placeholder="10" 
            required 
            id="trade-shares-input"
          />

          <!-- Quick Cash Preset Sizers -->
          <div class="preset-buttons">
            <button type="button" class="preset-btn" @click="setPresetAmount(0.25)">25%</button>
            <button type="button" class="preset-btn" @click="setPresetAmount(0.50)">50%</button>
            <button type="button" class="preset-btn" @click="setPresetAmount(0.75)">75%</button>
            <button type="button" class="preset-btn" @click="setPresetAmount(1.00)">Max</button>
          </div>
        </div>

        <!-- Live Trade Impact Preview Box -->
        <div class="impact-preview-box">
          <div class="impact-title">⚡ Trade Impact Preview</div>
          <div class="impact-row">
            <span>Estimated {{ tradeType === 'BUY' ? 'Total Cost' : 'Gross Proceeds' }}:</span>
            <strong class="text-white font-heading">${{ estimatedTotal.toFixed(2) }}</strong>
          </div>
          <div class="impact-row sub">
            <span>Remaining Cash Buffer:</span>
            <span :class="resultingCash >= 0 ? 'text-success' : 'text-danger'">
              ${{ resultingCash.toFixed(2) }}
            </span>
          </div>
          <div class="impact-row sub">
            <span>Estimated Asset Allocation:</span>
            <span>~{{ projectedAllocPct.toFixed(1) }}% of Portfolio</span>
          </div>
        </div>

        <!-- Error Message Banner -->
        <div v-if="error || isOverBudget" class="error-banner">
          {{ error || 'Insufficient cash balance to complete this buy order.' }}
        </div>

        <!-- Confirm CTA -->
        <button 
          type="submit" 
          :class="['btn', 'btn-block', 'btn-lg', tradeType === 'BUY' ? 'btn-success' : 'btn-danger']"
          :disabled="loading || isOverBudget || (tradeType === 'SELL' && shares > ownedShares)"
          id="confirm-trade-btn"
        >
          <span v-if="loading">Processing Order...</span>
          <span v-else>Confirm {{ tradeType }} Order (${{ estimatedTotal.toFixed(2) }})</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  tickers: { type: Array, default: () => ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'NVDA', 'SPY'] },
  priceMap: { type: Object, default: () => ({}) },
  initialTicker: { type: String, default: 'AAPL' },
  initialType: { type: String, default: 'BUY' },
  cashBalance: { type: Number, default: 100000 },
  holdings: { type: Array, default: () => [] },
  loading: Boolean,
  error: String
})

const emit = defineEmits(['close', 'execute-trade'])

const tradeType = ref(props.initialType)
const selectedTicker = ref(props.initialTicker)
const shares = ref(1)

const companyNames = {
  AAPL: 'Apple Inc.',
  MSFT: 'Microsoft Corp.',
  GOOGL: 'Alphabet Inc.',
  AMZN: 'Amazon.com Inc.',
  NVDA: 'NVIDIA Corp.',
  SPY: 'S&P 500 ETF Trust'
}

function getCompanyName(ticker) {
  return companyNames[ticker] || 'Equity'
}

watch(() => props.initialTicker, (val) => { if (val) selectedTicker.value = val })
watch(() => props.initialType, (val) => { if (val) tradeType.value = val })

const currentPrice = computed(() => {
  const pData = props.priceMap[selectedTicker.value]
  return pData && pData.price ? pData.price : 150.0
})

const estimatedTotal = computed(() => {
  return (shares.value || 0) * currentPrice.value
})

const ownedShares = computed(() => {
  const item = (props.holdings || []).find(h => h.ticker === selectedTicker.value)
  return item ? item.shares : 0
})

const isOverBudget = computed(() => {
  if (tradeType.value === 'BUY') {
    return estimatedTotal.value > props.cashBalance
  }
  return false
})

const resultingCash = computed(() => {
  if (tradeType.value === 'BUY') {
    return Math.max(0, props.cashBalance - estimatedTotal.value)
  }
  return props.cashBalance + estimatedTotal.value
})

const projectedAllocPct = computed(() => {
  const price = currentPrice.value
  const totalVal = props.cashBalance + (props.holdings || []).reduce((acc, h) => acc + (h.market_value || 0), 0)
  if (totalVal <= 0) return 0

  const currentHoldingVal = ownedShares.value * price
  const changeVal = tradeType.value === 'BUY' ? estimatedTotal.value : -estimatedTotal.value
  const finalVal = Math.max(0, currentHoldingVal + changeVal)

  return (finalVal / totalVal) * 100
})

function setPresetAmount(ratio) {
  const price = currentPrice.value
  if (price <= 0) return

  if (tradeType.value === 'BUY') {
    const budget = props.cashBalance * ratio
    const calculatedShares = Math.floor((budget / price) * 100) / 100
    shares.value = Math.max(0.01, calculatedShares)
  } else {
    const available = ownedShares.value * ratio
    shares.value = Math.max(0.01, Math.round(available * 100) / 100)
  }
}

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
  max-width: 460px;
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

.trade-type-toggle {
  display: flex;
  background: rgba(0, 0, 0, 0.4);
  padding: 4px;
  border-radius: 12px;
  margin-bottom: 1.25rem;
}

.toggle-btn {
  flex: 1;
  padding: 0.65rem;
  border: none;
  background: none;
  color: var(--text-muted);
  font-weight: 800;
  font-size: 0.85rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.04em;
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

.price-quote-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  padding: 0.75rem 1rem;
  border-radius: 12px;
  margin-bottom: 1.15rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quote-info {
  display: flex;
  flex-direction: column;
}

.quote-label {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.quote-price {
  font-size: 1.35rem;
  font-weight: 800;
  color: #a5b4fc;
}

.live-dot-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.7rem;
  color: #34d399;
  background: rgba(16, 185, 129, 0.15);
  padding: 0.2rem 0.5rem;
  border-radius: 9999px;
}

.pulse-dot {
  width: 6px;
  height: 6px;
  background-color: #10b981;
  border-radius: 50%;
}

.form-group {
  margin-bottom: 1.15rem;
}

.label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.35rem;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
}

.owned-shares-hint {
  font-size: 0.75rem;
  color: #a5b4fc;
  font-weight: 600;
}

.shares-input {
  font-size: 1.15rem;
  font-weight: 700;
}

.preset-buttons {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.preset-btn {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  color: var(--text-main);
  border-radius: 8px;
  padding: 0.35rem 0;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.preset-btn:hover {
  background: rgba(99, 102, 241, 0.2);
  border-color: var(--primary);
}

.impact-preview-box {
  background: rgba(0, 0, 0, 0.35);
  border: 1px dashed var(--border-color);
  padding: 0.85rem 1rem;
  border-radius: 12px;
  margin-bottom: 1.25rem;
}

.impact-title {
  font-size: 0.75rem;
  font-weight: 700;
  color: #a5b4fc;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
}

.impact-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.impact-row.sub {
  font-size: 0.75rem;
  color: var(--text-muted);
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
