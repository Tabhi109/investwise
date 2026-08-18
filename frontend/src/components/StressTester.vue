<template>
  <div class="glass-card stress-tester-card">
    <div class="tester-header">
      <div class="title-wrap">
        <span class="chip-accent">SCENARIO LAB</span>
        <h3 class="card-title">Portfolio Stress Tester</h3>
      </div>
      <span class="badge badge-primary">Interactive Shock Simulator</span>
    </div>

    <p class="tester-desc">
      Simulate real-world macroeconomic market events and immediately see how your current portfolio holdings, beta, and cash buffer will react.
    </p>

    <!-- Preset Scenarios -->
    <div class="preset-grid">
      <button 
        v-for="scenario in scenarios" 
        :key="scenario.id"
        class="scenario-btn"
        :class="{ 'active': selectedScenario.id === scenario.id }"
        @click="selectScenario(scenario)"
      >
        <span class="scenario-icon">{{ scenario.icon }}</span>
        <span class="scenario-name">{{ scenario.name }}</span>
        <span class="scenario-shock" :class="scenario.shock >= 0 ? 'text-success' : 'text-danger'">
          {{ scenario.shock >= 0 ? '+' : '' }}{{ scenario.shock }}%
        </span>
      </button>
    </div>

    <!-- Custom Slider -->
    <div class="slider-box">
      <div class="slider-header">
        <span class="slider-label">Custom Market Shock:</span>
        <span class="slider-val font-heading" :class="customShock >= 0 ? 'text-success' : 'text-danger'">
          {{ customShock >= 0 ? '+' : '' }}{{ customShock }}%
        </span>
      </div>
      <input 
        type="range" 
        min="-30" 
        max="30" 
        step="1" 
        v-model.number="customShock" 
        class="shock-slider"
        @input="onSliderInput"
      />
      <div class="slider-ticks">
        <span>-30% (Severe Crash)</span>
        <span>0% (Neutral)</span>
        <span>+30% (Mega Rally)</span>
      </div>
    </div>

    <!-- Projected Impact Result Card -->
    <div class="projection-card" :class="projectedPnl >= 0 ? 'bull-card' : 'bear-card'">
      <div class="projection-grid">
        <div class="proj-item">
          <span class="proj-label">Estimated P&L Impact</span>
          <span class="proj-val font-heading" :class="projectedPnl >= 0 ? 'text-success' : 'text-danger'">
            {{ projectedPnl >= 0 ? '+' : '' }}${{ formatMoney(projectedPnl) }}
          </span>
          <span class="proj-pct">
            ({{ projectedPnlPct >= 0 ? '+' : '' }}{{ projectedPnlPct.toFixed(2) }}% of Portfolio)
          </span>
        </div>

        <div class="proj-item">
          <span class="proj-label">Projected Net Worth</span>
          <span class="proj-val font-heading text-white">
            ${{ formatMoney(projectedNetWorth) }}
          </span>
          <span class="proj-subtext">Current: ${{ formatMoney(totalPortfolioValue) }}</span>
        </div>
      </div>

      <!-- Educational Takeaway -->
      <div class="educational-takeaway">
        <span class="takeaway-icon">💡</span>
        <p class="takeaway-text">{{ projectedTakeaway }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  summary: {
    type: Object,
    default: () => ({
      total_value: 100000,
      cash_balance: 100000,
      holdings_value: 0
    })
  },
  beta: {
    type: Number,
    default: 1.0
  }
})

const scenarios = [
  { id: 'crash', name: 'Market Flash Crash', icon: '⚡', shock: -15, desc: 'Sharp broader market liquidation like March 2020.' },
  { id: 'ratehike', name: 'Rate Hike Scare', icon: '🏦', shock: -8, desc: 'Central bank tightening interest rates dampening growth.' },
  { id: 'neutral', name: 'Sideways Flat', icon: '⚖️', shock: 0, desc: 'Market trades in a tight range without clear trend.' },
  { id: 'techrally', name: 'AI & Tech Surge', icon: '🚀', shock: 12, desc: 'Broad tech momentum driving market multiples higher.' },
  { id: 'bullrun', name: 'Mega Bull Wave', icon: '🐂', shock: 25, desc: 'Euphoric market rally lifting nearly all equity assets.' }
]

const selectedScenario = ref(scenarios[0])
const customShock = ref(-15)

function selectScenario(scenario) {
  selectedScenario.value = scenario
  customShock.value = scenario.shock
}

function onSliderInput() {
  selectedScenario.value = { id: 'custom', name: 'Custom Slider', icon: '🎛️', shock: customShock.value }
}

const totalPortfolioValue = computed(() => props.summary.total_value || 100000)
const holdingsValue = computed(() => props.summary.holdings_value || 0)
const effectiveBeta = computed(() => props.beta || 1.0)

const projectedPnl = computed(() => {
  // Dollar shock based on stock holdings value multiplied by effective Beta and shock %
  const marketShockRatio = customShock.value / 100
  return holdingsValue.value * effectiveBeta.value * marketShockRatio
})

const projectedPnlPct = computed(() => {
  if (totalPortfolioValue.value <= 0) return 0
  return (projectedPnl.value / totalPortfolioValue.value) * 100
})

const projectedNetWorth = computed(() => {
  return Math.max(0, totalPortfolioValue.value + projectedPnl.value)
})

const projectedTakeaway = computed(() => {
  const cashPct = totalPortfolioValue.value > 0 ? (props.summary.cash_balance / totalPortfolioValue.value) * 100 : 100
  
  if (holdingsValue.value === 0) {
    return 'Your portfolio is currently 100% in cash. You are fully insulated from market downturns, but cash yields 0% growth against inflation!'
  }

  if (customShock.value < 0) {
    if (cashPct > 30) {
      return `Because you hold ${cashPct.toFixed(0)}% in cash, your actual portfolio drop is cushioned to only ${Math.abs(projectedPnlPct.value).toFixed(1)}%, giving you dry powder to buy discounted quality stocks!`
    }
    return `With ${cashPct.toFixed(0)}% cash and a Beta of ${effectiveBeta.value.toFixed(2)}, your equity portfolio absorbs the drop directly. Consider holding a small cash cushion to rebalance during corrections.`
  } else if (customShock.value > 0) {
    return `In this rally, your stock holdings generate +$${formatMoney(projectedPnl.value)}. Having a Beta of ${effectiveBeta.value.toFixed(2)} amplifies your upside participation!`
  }
  return 'During flat markets, dividend yield, risk mitigation, and selective stock picking drive outperformance.'
})

function formatMoney(val) {
  if (val === undefined || val === null) return '0.00'
  return Math.abs(val).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
</script>

<style scoped>
.stress-tester-card {
  padding: 1.5rem;
  margin-bottom: 1.75rem;
  background: rgba(17, 24, 39, 0.85);
}

.tester-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.title-wrap {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.chip-accent {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.25), rgba(6, 182, 212, 0.2));
  color: #c4b5fd;
  border: 1px solid rgba(139, 92, 246, 0.4);
  font-size: 0.65rem;
  font-weight: 800;
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
  letter-spacing: 0.08em;
}

.card-title {
  font-size: 1.15rem;
  font-weight: 700;
}

.tester-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  line-height: 1.45;
  margin-bottom: 1.25rem;
}

.preset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
  gap: 0.65rem;
  margin-bottom: 1.25rem;
}

@media (max-width: 640px) {
  .preset-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.scenario-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 0.75rem 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.scenario-btn:hover {
  background: rgba(255, 255, 255, 0.07);
  border-color: var(--border-active);
}

.scenario-btn.active {
  background: rgba(99, 102, 241, 0.18);
  border-color: #818cf8;
  box-shadow: 0 0 16px rgba(99, 102, 241, 0.25);
}

.scenario-icon {
  font-size: 1.3rem;
  margin-bottom: 0.25rem;
}

.scenario-name {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-main);
  margin-bottom: 0.2rem;
}

.scenario-shock {
  font-size: 0.85rem;
  font-weight: 800;
  font-family: var(--font-heading);
}

.slider-box {
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1.25rem;
}

.slider-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.slider-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
}

.slider-val {
  font-size: 1.15rem;
  font-weight: 800;
}

.shock-slider {
  width: 100%;
  accent-color: var(--primary);
  cursor: pointer;
  margin-bottom: 0.35rem;
}

.slider-ticks {
  display: flex;
  justify-content: space-between;
  font-size: 0.68rem;
  color: var(--text-dim);
}

.projection-card {
  border-radius: 14px;
  padding: 1.25rem;
  border: 1px solid var(--border-color);
  transition: all 0.3s;
}

.bear-card {
  background: linear-gradient(135deg, rgba(244, 63, 94, 0.08), rgba(15, 23, 42, 0.6));
  border-color: rgba(244, 63, 94, 0.3);
}

.bull-card {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.08), rgba(15, 23, 42, 0.6));
  border-color: rgba(16, 185, 129, 0.3);
}

.projection-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

@media (max-width: 500px) {
  .projection-grid {
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }
}

.proj-item {
  display: flex;
  flex-direction: column;
}

.proj-label {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  font-weight: 600;
}

.proj-val {
  font-size: 1.5rem;
  font-weight: 800;
  margin: 0.2rem 0;
}

.proj-pct {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.proj-subtext {
  font-size: 0.75rem;
  color: var(--text-dim);
}

.text-white { color: white; }
.text-success { color: var(--success); }
.text-danger { color: var(--danger); }

.educational-takeaway {
  display: flex;
  gap: 0.75rem;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 10px;
  padding: 0.85rem 1rem;
  font-size: 0.825rem;
  color: #e2e8f0;
  line-height: 1.45;
}

.takeaway-icon {
  font-size: 1.25rem;
}
</style>
