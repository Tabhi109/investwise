<template>
  <div class="glass-card history-card">
    <div class="card-header">
      <h3 class="card-title">Order Audit Log</h3>
    </div>

    <div class="table-responsive">
      <table class="data-table">
        <thead>
          <tr>
            <th>Order ID</th>
            <th>Type</th>
            <th>Asset Ticker</th>
            <th>Shares</th>
            <th>Executed Price</th>
            <th>Total Value</th>
            <th>Timestamp</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!transactions || transactions.length === 0">
            <td colspan="7" class="empty-cell">No trade order history recorded.</td>
          </tr>
          <tr v-for="t in transactions" :key="t.id">
            <td class="font-mono">#{{ t.id }}</td>
            <td>
              <span :class="['type-badge', t.transaction_type === 'BUY' ? 'type-buy' : 'type-sell']">
                {{ t.transaction_type }}
              </span>
            </td>
            <td class="font-bold">{{ t.ticker }}</td>
            <td>{{ formatNum(t.shares, 4) }}</td>
            <td>${{ formatNum(t.price, 2) }}</td>
            <td class="font-bold">${{ formatNum(t.shares * t.price, 2) }}</td>
            <td class="time-cell">{{ formatDate(t.transaction_time) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
defineProps({
  transactions: Array
})

function formatNum(val, decimals = 2) {
  if (val === undefined || val === null) return '0.00'
  return val.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
}

function formatDate(isoStr) {
  if (!isoStr) return '--'
  const d = new Date(isoStr)
  return d.toLocaleString()
}
</script>

<style scoped>
.history-card {
  margin-bottom: 2rem;
  padding: 1.25rem;
}

.card-header {
  margin-bottom: 1rem;
}

.card-title {
  font-size: 1.15rem;
  font-weight: 700;
}

.table-responsive {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
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

.font-mono {
  font-family: monospace;
  color: var(--text-dim);
}

.font-bold {
  font-weight: 700;
}

.type-badge {
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
  font-size: 0.725rem;
  font-weight: 800;
}

.type-buy {
  background: var(--success-glow);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.type-sell {
  background: var(--danger-glow);
  color: #fb7185;
  border: 1px solid rgba(244, 63, 94, 0.3);
}

.time-cell {
  color: var(--text-muted);
  font-size: 0.8rem;
}

.empty-cell {
  text-align: center;
  color: var(--text-dim);
  padding: 2rem !important;
}
</style>
