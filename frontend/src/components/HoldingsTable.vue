<template>
  <div class="glass-card table-card">
    <div class="card-header">
      <div class="header-left">
        <h3 class="card-title">Portfolio Holdings</h3>
        <span v-if="holdings && holdings.length > 0" class="badge badge-neutral">
          {{ holdings.length }} {{ holdings.length === 1 ? 'Position' : 'Positions' }}
        </span>
      </div>
      <button @click="$emit('open-trade', 'BUY')" class="btn btn-primary btn-sm" id="add-holding-btn">
        + Trade Stock
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="!holdings || holdings.length === 0" class="empty-state">
      <div class="empty-icon">📊</div>
      <h4>No Active Stock Holdings</h4>
      <p>Your demo capital ($100,000) is currently 100% in cash. Purchase your first shares to start building wealth and analyzing real-time risk metrics!</p>
      <button @click="$emit('open-trade', 'BUY')" class="btn btn-success btn-sm">
        Buy Your First Stock
      </button>
    </div>

    <template v-else>
      <!-- Mobile Card List View (< 768px) -->
      <div class="mobile-holdings-list">
        <div 
          v-for="h in holdings" 
          :key="'mob-' + h.ticker"
          class="holding-mobile-card"
        >
          <div class="holding-card-top">
            <div class="holding-main">
              <span class="symbol-badge">{{ h.ticker }}</span>
              <div class="holding-meta">
                <span class="holding-shares">{{ formatNum(h.shares, 2) }} shares</span>
                <span class="holding-avg">Avg: ${{ formatNum(h.average_buy_price, 2) }}</span>
              </div>
            </div>
            <div class="holding-value-block">
              <span class="holding-mkt-val font-heading">${{ formatNum(h.market_value, 2) }}</span>
              <span :class="['holding-pnl-chip', h.unrealized_pnl >= 0 ? 'text-success' : 'text-danger']">
                {{ h.unrealized_pnl >= 0 ? '+' : '' }}${{ formatNum(h.unrealized_pnl, 2) }}
              </span>
            </div>
          </div>

          <!-- Allocation Bar -->
          <div class="card-alloc-row">
            <div class="alloc-bar-track">
              <div class="alloc-bar" :style="{ width: Math.min(100, (h.allocation_pct * 100)) + '%' }"></div>
            </div>
            <span class="alloc-pct-text">{{ (h.allocation_pct * 100).toFixed(1) }}% of Portfolio</span>
          </div>

          <!-- Quick Action Buttons -->
          <div class="card-actions-row">
            <button @click="$emit('trade-item', { ticker: h.ticker, type: 'BUY' })" class="btn btn-success btn-xs">
              + Buy More
            </button>
            <button @click="$emit('trade-item', { ticker: h.ticker, type: 'SELL' })" class="btn btn-danger btn-xs">
              - Sell Position
            </button>
          </div>
        </div>
      </div>

      <!-- Desktop Table View (>= 768px) -->
      <div class="table-responsive desktop-only-table">
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
            <tr v-for="h in holdings" :key="'desk-' + h.ticker">
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
                  <div class="alloc-bar" :style="{ width: Math.min(100, (h.allocation_pct * 100)) + '%' }"></div>
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
    </template>
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
  margin-bottom: 1.5rem;
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

.btn-sm {
  padding: 0.4rem 0.85rem;
  font-size: 0.8rem;
}

.btn-xs {
  padding: 0.35rem 0.75rem;
  font-size: 0.75rem;
  border-radius: 8px;
  flex: 1;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 2.5rem 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.empty-state h4 {
  font-size: 1.1rem;
  margin-bottom: 0.4rem;
}

.empty-state p {
  color: var(--text-muted);
  font-size: 0.85rem;
  max-width: 450px;
  line-height: 1.45;
  margin-bottom: 1.25rem;
}

/* Mobile Holdings Card List */
.mobile-holdings-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

@media (min-width: 768px) {
  .mobile-holdings-list {
    display: none;
  }
}

@media (max-width: 767px) {
  .desktop-only-table {
    display: none;
  }
}

.holding-mobile-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0.85rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.holding-card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.holding-main {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.holding-meta {
  display: flex;
  flex-direction: column;
}

.holding-shares {
  font-size: 0.85rem;
  font-weight: 700;
}

.holding-avg {
  font-size: 0.7rem;
  color: var(--text-dim);
}

.holding-value-block {
  text-align: right;
  display: flex;
  flex-direction: column;
}

.holding-mkt-val {
  font-size: 1.05rem;
  font-weight: 800;
}

.holding-pnl-chip {
  font-size: 0.75rem;
  font-weight: 700;
}

.card-alloc-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.alloc-bar-track {
  flex: 1;
  height: 5px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  overflow: hidden;
}

.alloc-pct-text {
  font-size: 0.7rem;
  color: var(--text-dim);
  min-width: 90px;
  text-align: right;
}

.card-actions-row {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.2rem;
}

/* Desktop Data Table */
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
  font-weight: 800;
  font-size: 0.85rem;
}

.font-bold { font-weight: 700; }
.text-success { color: var(--success); font-weight: 600; }
.text-danger { color: var(--danger); font-weight: 600; }

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
