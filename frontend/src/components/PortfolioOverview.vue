<template>
  <div class="overview-grid">
    <div class="glass-card metric-card hero-card">
      <div class="metric-label">Total Portfolio Net Worth</div>
      <div class="metric-value font-heading">${{ formatMoney(summary.total_value) }}</div>
      <div class="metric-footer">
        <span class="badge badge-neutral">Cash: ${{ formatMoney(summary.cash_balance) }}</span>
        <span class="badge badge-neutral">Holdings: ${{ formatMoney(summary.holdings_value) }}</span>
      </div>
    </div>

    <div class="glass-card metric-card">
      <div class="metric-label">Unrealized P&L</div>
      <div :class="['metric-value', summary.unrealized_pnl >= 0 ? 'text-success' : 'text-danger']">
        {{ summary.unrealized_pnl >= 0 ? '+' : '' }}${{ formatMoney(summary.unrealized_pnl) }}
      </div>
      <div class="metric-subtext">Open positions valuation gain/loss</div>
    </div>

    <div class="glass-card metric-card">
      <div class="metric-label">Realized P&L</div>
      <div :class="['metric-value', summary.realized_pnl >= 0 ? 'text-success' : 'text-danger']">
        {{ summary.realized_pnl >= 0 ? '+' : '' }}${{ formatMoney(summary.realized_pnl) }}
      </div>
      <div class="metric-subtext">Closed trade execution gain/loss</div>
    </div>

    <div class="glass-card metric-card action-card">
      <div class="metric-label">Quick Execution</div>
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
</template>

<script setup>
defineProps({
  summary: {
    type: Object,
    default: () => ({
      total_value: 0,
      cash_balance: 0,
      holdings_value: 0,
      unrealized_pnl: 0,
      realized_pnl: 0
    })
  }
})

defineEmits(['open-trade'])

function formatMoney(val) {
  if (val === undefined || val === null) return '0.00'
  return val.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
</script>

<style scoped>
.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.25rem;
  margin-bottom: 1.75rem;
}

.metric-card {
  padding: 1.35rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.hero-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.1));
  border-color: rgba(99, 102, 241, 0.3);
}

.metric-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.4rem;
}

.metric-value {
  font-size: 1.8rem;
  font-weight: 800;
  line-height: 1.2;
}

.font-heading {
  font-family: var(--font-heading);
}

.text-success {
  color: var(--success);
}

.text-danger {
  color: var(--danger);
}

.metric-footer {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.metric-subtext {
  font-size: 0.75rem;
  color: var(--text-dim);
  margin-top: 0.4rem;
}

.action-buttons {
  display: flex;
  gap: 0.65rem;
  margin-top: 0.75rem;
}

.action-buttons button {
  flex: 1;
}
</style>
