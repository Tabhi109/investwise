<template>
  <div class="glass-card table-card">
    <div class="card-header">
      <h3 class="card-title">Portfolio Holdings</h3>
      <button @click="$emit('open-trade', 'BUY')" class="btn btn-primary btn-sm" id="add-holding-btn">
        + Trade Stock
      </button>
    </div>

    <div class="table-responsive">
      <table class="data-table">
        <thead>
          <tr>
            <th>Asset</th>
            <th>Shares</th>
            <th>Avg Buy Price</th>
            <th>Current Price</th>
            <th>Market Value</th>
            <th>Unrealized P&L</th>
            <th>Allocation %</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!holdings || holdings.length === 0">
            <td colspan="8" class="empty-cell">
              No stock holdings owned yet. Use "+ Trade Stock" to buy your first position!
            </td>
          </tr>
          <tr v-for="h in holdings" :key="h.ticker">
            <td class="symbol-cell">
              <span class="symbol-badge">{{ h.ticker }}</span>
            </td>
            <td>{{ formatNum(h.shares, 4) }}</td>
            <td>${{ formatNum(h.average_buy_price, 2) }}</td>
            <td class="font-bold">${{ formatNum(h.current_price, 2) }}</td>
            <td class="font-bold">${{ formatNum(h.market_value, 2) }}</td>
            <td :class="h.unrealized_pnl >= 0 ? 'text-success' : 'text-danger'">
              {{ h.unrealized_pnl >= 0 ? '+' : '' }}${{ formatNum(h.unrealized_pnl, 2) }}
            </td>
            <td>
              <div class="alloc-bar-container">
                <div class="alloc-bar" :style="{ width: (h.allocation_pct * 100) + '%' }"></div>
                <span>{{ (h.allocation_pct * 100).toFixed(1) }}%</span>
              </div>
            </td>
            <td>
              <div class="action-cell">
                <button @click="$emit('trade-item', { ticker: h.ticker, type: 'BUY' })" class="btn btn-success btn-xs">Buy</button>
                <button @click="$emit('trade-item', { ticker: h.ticker, type: 'SELL' })" class="btn btn-danger btn-xs">Sell</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
defineProps({
  holdings: Array
})

defineEmits(['open-trade', 'trade-item'])

function formatNum(val, decimals = 2) {
  if (val === undefined || val === null) return '0.00'
  return val.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
}
</script>

<style scoped>
.table-card {
  margin-bottom: 1.75rem;
  padding: 1.25rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.card-title {
  font-size: 1.15rem;
  font-weight: 700;
}

.btn-sm {
  padding: 0.4rem 0.85rem;
  font-size: 0.8rem;
}

.btn-xs {
  padding: 0.25rem 0.5rem;
  font-size: 0.725rem;
  border-radius: 6px;
}

.table-responsive {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}

.data-table th {
  color: var(--text-muted);
  font-weight: 600;
  padding: 0.75rem 0.85rem;
  border-bottom: 1px solid var(--border-color);
  font-size: 0.75rem;
  text-transform: uppercase;
}

.data-table td {
  padding: 0.85rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.data-table tr:last-child td {
  border-bottom: none;
}

.symbol-badge {
  background: rgba(99, 102, 241, 0.15);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.3);
  padding: 0.25rem 0.6rem;
  border-radius: 6px;
  font-weight: 700;
  font-size: 0.8rem;
}

.font-bold {
  font-weight: 700;
}

.text-success {
  color: var(--success);
  font-weight: 600;
}

.text-danger {
  color: var(--danger);
  font-weight: 600;
}

.empty-cell {
  text-align: center;
  color: var(--text-dim);
  padding: 2rem !important;
}

.action-cell {
  display: flex;
  gap: 0.35rem;
}

.alloc-bar-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 100px;
}

.alloc-bar {
  height: 6px;
  background: linear-gradient(90deg, var(--primary), var(--accent));
  border-radius: 3px;
  max-width: 70px;
}
</style>
