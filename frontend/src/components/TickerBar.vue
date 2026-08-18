<template>
  <div class="ticker-bar glass-card">
    <div class="ticker-scroll-container">
      <div 
        v-for="(data, ticker) in priceMap" 
        :key="ticker" 
        class="ticker-card"
        :class="{ 'tick-up': data.flash === 'up', 'tick-down': data.flash === 'down' }"
        @click="$emit('select-ticker', ticker)"
        :title="`Trade ${ticker}`"
      >
        <div class="ticker-header">
          <div class="symbol-wrap">
            <span class="ticker-symbol">{{ ticker }}</span>
            <span class="ticker-name-sub">{{ getCompanyName(ticker) }}</span>
          </div>
          <span :class="['pct-badge', data.change >= 0 ? 'pct-up' : 'pct-down']">
            {{ data.change >= 0 ? '+' : '' }}{{ data.change ? data.change.toFixed(2) : '0.00' }}%
          </span>
        </div>

        <div class="ticker-footer">
          <span class="ticker-price">
            ${{ data.price ? data.price.toFixed(2) : '--.--' }}
          </span>
          <span class="trade-cta-chip">Trade &rsaquo;</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  priceMap: Object
})
defineEmits(['select-ticker'])

const companyNames = {
  AAPL: 'Apple',
  MSFT: 'Microsoft',
  GOOGL: 'Alphabet',
  AMZN: 'Amazon',
  NVDA: 'Nvidia',
  SPY: 'S&P 500 ETF'
}

function getCompanyName(ticker) {
  return companyNames[ticker] || 'Stock'
}
</script>

<style scoped>
.ticker-bar {
  margin-bottom: 1.25rem;
  padding: 0.65rem 0.85rem;
  overflow: hidden;
}

.ticker-scroll-container {
  display: flex;
  gap: 0.75rem;
  overflow-x: auto;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;
  padding-bottom: 0.25rem;
}

/* Hide scrollbar on touch devices for app-like feel */
.ticker-scroll-container::-webkit-scrollbar {
  height: 4px;
}

.ticker-card {
  flex: 0 0 170px;
  scroll-snap-align: start;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0.65rem 0.85rem;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

@media (min-width: 1100px) {
  .ticker-scroll-container {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    overflow-x: visible;
  }

  .ticker-card {
    flex: auto;
  }
}

.ticker-card:hover {
  transform: translateY(-2px);
  border-color: var(--primary);
  background: rgba(255, 255, 255, 0.06);
}

.ticker-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.35rem;
}

.symbol-wrap {
  display: flex;
  flex-direction: column;
}

.ticker-symbol {
  font-weight: 800;
  font-size: 0.9rem;
  color: var(--text-main);
}

.ticker-name-sub {
  font-size: 0.65rem;
  color: var(--text-dim);
}

.pct-badge {
  font-size: 0.725rem;
  font-weight: 700;
  padding: 0.15rem 0.4rem;
  border-radius: 6px;
}

.pct-up {
  color: #34d399;
  background: rgba(16, 185, 129, 0.15);
}

.pct-down {
  color: #fb7185;
  background: rgba(244, 63, 94, 0.15);
}

.ticker-footer {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.ticker-price {
  font-family: var(--font-heading);
  font-size: 1.15rem;
  font-weight: 800;
}

.trade-cta-chip {
  font-size: 0.65rem;
  font-weight: 600;
  color: #a5b4fc;
}
</style>
