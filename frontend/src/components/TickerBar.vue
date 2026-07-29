<template>
  <div class="ticker-bar glass-card">
    <div class="ticker-scroll">
      <div 
        v-for="(data, ticker) in priceMap" 
        :key="ticker" 
        class="ticker-card"
        :class="{ 'tick-up': data.flash === 'up', 'tick-down': data.flash === 'down' }"
        @click="$emit('select-ticker', ticker)"
      >
        <div class="ticker-header">
          <span class="ticker-symbol">{{ ticker }}</span>
          <span :class="data.change >= 0 ? 'pct-up' : 'pct-down'">
            {{ data.change >= 0 ? '+' : '' }}{{ data.change ? data.change.toFixed(2) : '0.00' }}%
          </span>
        </div>
        <div class="ticker-price">
          ${{ data.price ? data.price.toFixed(2) : '--.--' }}
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
</script>

<style scoped>
.ticker-bar {
  margin-bottom: 1.5rem;
  padding: 0.75rem 1rem;
  overflow: hidden;
}

.ticker-scroll {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 0.75rem;
}

.ticker-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0.65rem 0.85rem;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s;
}

.ticker-card:hover {
  transform: translateY(-2px);
  border-color: var(--primary);
  background: rgba(255, 255, 255, 0.06);
}

.ticker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.25rem;
}

.ticker-symbol {
  font-weight: 700;
  font-size: 0.85rem;
  color: var(--text-main);
}

.pct-up {
  color: var(--success);
  font-size: 0.75rem;
  font-weight: 600;
}

.pct-down {
  color: var(--danger);
  font-size: 0.75rem;
  font-weight: 600;
}

.ticker-price {
  font-family: var(--font-heading);
  font-size: 1.1rem;
  font-weight: 700;
}
</style>
