<template>
  <div class="glass-card history-card">
    <div class="card-header">
      <div class="header-left">
        <h3 class="card-title">Order Audit Log</h3>
        <span v-if="transactions && transactions.length > 0" class="badge badge-neutral">
          {{ transactions.length }} {{ transactions.length === 1 ? 'Order' : 'Orders' }}
        </span>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!transactions || transactions.length === 0" class="empty-state">
      <div class="empty-icon">📝</div>
      <p>No trade orders recorded yet. Executed market orders will appear here for full audit traceability.</p>
    </div>

    <template v-else>
      <!-- Mobile Order List View (< 768px) -->
      <div class="mobile-order-list">
        <div 
          v-for="t in transactions" 
          :key="'mob-' + t.id"
          class="order-mobile-item"
        >
          <div class="order-item-left">
            <span :class="['type-badge', t.transaction_type === 'BUY' ? 'type-buy' : 'type-sell']">
              {{ t.transaction_type }}
            </span>
            <div class="order-details">
              <span class="order-ticker">{{ t.ticker }}</span>
              <span class="order-time">{{ formatDate(t.transaction_time) }}</span>
            </div>
          </div>
          <div class="order-item-right">
            <span class="order-val font-heading">${{ formatNum(t.shares * t.price, 2) }}</span>
            <span class="order-sub">{{ formatNum(t.shares, 2) }} shares @ ${{ formatNum(t.price, 2) }}</span>
          </div>
        </div>
      </div>

      <!-- Desktop Table View (>= 768px) -->
      <div class="table-responsive desktop-only-table">
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
            <tr v-for="t in transactions" :key="'desk-' + t.id">
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
    </template>
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
  return d.toLocaleDateString() + ' ' + d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.history-card {
  margin-bottom: 2rem;
  padding: 1.25rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.card-title {
  font-size: 1.15rem;
  font-weight: 700;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.empty-icon {
  font-size: 2rem;
  margin-bottom: 0.35rem;
}

/* Mobile Order List */
.mobile-order-list {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

@media (min-width: 768px) {
  .mobile-order-list {
    display: none;
  }
}

@media (max-width: 767px) {
  .desktop-only-table {
    display: none;
  }
}

.order-mobile-item {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0.75rem 0.9rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-item-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.order-details {
  display: flex;
  flex-direction: column;
}

.order-ticker {
  font-weight: 700;
  font-size: 0.9rem;
}

.order-time {
  font-size: 0.7rem;
  color: var(--text-dim);
}

.order-item-right {
  text-align: right;
  display: flex;
  flex-direction: column;
}

.order-val {
  font-size: 0.95rem;
  font-weight: 700;
}

.order-sub {
  font-size: 0.7rem;
  color: var(--text-dim);
}

/* Desktop Data Table */
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

.data-table tr:last-child td {
  border-bottom: none;
}

.font-mono {
  font-family: monospace;
  color: var(--text-dim);
}

.font-bold { font-weight: 700; }

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
</style>
