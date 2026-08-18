<template>
  <div class="overview-section">
    <div class="overview-grid">
      <!-- Main Net Worth Hero Card -->
      <div class="glass-card metric-card hero-card">
        <div class="card-top-row">
          <span class="metric-label">Total Portfolio Net Worth</span>
          <button class="health-pill" @click="$emit('explain-metric', 'health')">
            <span class="health-dot" :class="healthScoreClass"></span>
            <span>Health: <strong>{{ healthScore }}/100</strong></span>
            <span class="info-icon">ℹ️</span>
          </button>
        </div>

        <div class="metric-value font-heading">${{ formatMoney(summary.total_value) }}</div>

        <!-- Asset Allocation Ratio Visual Bar -->
        <div class="allocation-bar-wrap">
          <div class="alloc-ratio-bar">
            <div 
              class="ratio-fill-equity" 
              :style="{ width: equityRatio + '%' }"
              title="Equities"
            ></div>
            <div 
              class="ratio-fill-cash" 
              :style="{ width: cashRatio + '%' }"
              title="Cash"
            ></div>
          </div>
          <div class="alloc-labels">
            <span class="alloc-lbl"><span class="dot equity-dot"></span> Stocks: ${{ formatMoney(summary.holdings_value) }} ({{ equityRatio.toFixed(0) }}%)</span>
            <span class="alloc-lbl"><span class="dot cash-dot"></span> Cash: ${{ formatMoney(summary.cash_balance) }} ({{ cashRatio.toFixed(0) }}%)</span>
          </div>
        </div>
      </div>

      <!-- Unrealized P&L Card -->
      <div class="glass-card metric-card">
        <div class="metric-label">Unrealized P&L (Open Positions)</div>
        <div :class="['metric-value', summary.unrealized_pnl >= 0 ? 'text-success' : 'text-danger']">
          {{ summary.unrealized_pnl >= 0 ? '+' : '' }}${{ formatMoney(summary.unrealized_pnl) }}
        </div>
        <div class="metric-subtext">
          <span>Live market price movement</span>
        </div>
      </div>

      <!-- Realized P&L Card -->
      <div class="glass-card metric-card">
        <div class="metric-label">Realized P&L (Closed Trades)</div>
        <div :class="['metric-value', summary.realized_pnl >= 0 ? 'text-success' : 'text-danger']">
          {{ summary.realized_pnl >= 0 ? '+' : '' }}${{ formatMoney(summary.realized_pnl) }}
        </div>
        <div class="metric-subtext">
          <span>Locked-in trading gains/losses</span>
        </div>
      </div>

      <!-- Quick Execution Action Card -->
      <div class="glass-card metric-card action-card">
        <div class="metric-label">Quick Action</div>
        <div class="action-buttons">
          <button @click="$emit('open-trade', 'BUY')" class="btn btn-success" id="quick-buy-btn">
            + Buy Stock
          </button>
          <button @click="$emit('open-trade', 'SELL')" class="btn btn-danger" id="quick-sell-btn">
            - Sell Stock
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  summary: {
    type: Object,
    default: () => ({
      total_value: 0,
      cash_balance: 0,
      holdings_value: 0,
      unrealized_pnl: 0,
      realized_pnl: 0,
      holdings: []
    })
  }
})

defineEmits(['open-trade', 'explain-metric'])

const totalVal = computed(() => props.summary.total_value || 1)
const equityRatio = computed(() => Math.min(100, Math.max(0, (props.summary.holdings_value / totalVal.value) * 100)))
const cashRatio = computed(() => Math.min(100, Math.max(0, (props.summary.cash_balance / totalVal.value) * 100)))

// Compute dynamic 0-100 portfolio health score
const healthScore = computed(() => {
  let score = 50
  const holdingsCount = (props.summary.holdings || []).length

  // Diversification points
  if (holdingsCount >= 4) score += 25
  else if (holdingsCount >= 2) score += 15
  else if (holdingsCount === 1) score += 5

  // Cash buffer points (having 5% - 40% cash is healthy)
  const cashPct = cashRatio.value
  if (cashPct >= 5 && cashPct <= 40) score += 15
  else if (cashPct > 40 && cashPct <= 80) score += 10

  // Profitability points
  if (props.summary.unrealized_pnl > 0) score += 10

  return Math.min(100, score)
})

const healthScoreClass = computed(() => {
  if (healthScore.value >= 80) return 'health-green'
  if (healthScore.value >= 60) return 'health-yellow'
  return 'health-red'
})

function formatMoney(val) {
  if (val === undefined || val === null) return '0.00'
  return val.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
</script>

<style scoped>
.overview-section {
  margin-bottom: 1.5rem;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

@media (min-width: 900px) {
  .overview-grid {
    grid-template-columns: 1.5fr 1fr 1fr 1.1fr;
  }
}

.metric-card {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border-radius: 14px;
}

.hero-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.1));
  border-color: rgba(99, 102, 241, 0.35);
}

.card-top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.35rem;
}

.metric-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.health-pill {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid var(--border-color);
  padding: 0.2rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.7rem;
  color: var(--text-main);
  cursor: pointer;
  transition: all 0.2s;
}

.health-pill:hover {
  border-color: var(--primary);
  background: rgba(99, 102, 241, 0.2);
}

.health-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.health-green { background: #10b981; box-shadow: 0 0 8px #10b981; }
.health-yellow { background: #f59e0b; box-shadow: 0 0 8px #f59e0b; }
.health-red { background: #f43f5e; box-shadow: 0 0 8px #f43f5e; }

.info-icon {
  font-size: 0.65rem;
  opacity: 0.7;
}

.metric-value {
  font-size: 1.75rem;
  font-weight: 800;
  line-height: 1.2;
  margin: 0.2rem 0;
}

.allocation-bar-wrap {
  margin-top: 0.65rem;
}

.alloc-ratio-bar {
  display: flex;
  height: 6px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 0.4rem;
}

.ratio-fill-equity {
  background: linear-gradient(90deg, #6366f1, #8b5cf6);
  transition: width 0.3s ease;
}

.ratio-fill-cash {
  background: #10b981;
  transition: width 0.3s ease;
}

.alloc-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.7rem;
  color: var(--text-muted);
}

.dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 0.2rem;
}

.equity-dot { background: #818cf8; }
.cash-dot { background: #10b981; }

.metric-subtext {
  font-size: 0.75rem;
  color: var(--text-dim);
}

.text-success { color: var(--success); }
.text-danger { color: var(--danger); }

.action-buttons {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.action-buttons button {
  flex: 1;
}
</style>
