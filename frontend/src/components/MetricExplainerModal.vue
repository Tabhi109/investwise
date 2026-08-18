<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content glass-card modal-animate-sheet">
      <div class="modal-header">
        <div class="header-left">
          <div class="metric-icon-badge">💡</div>
          <div>
            <span class="chip-sm">FINANCIAL CONCEPT</span>
            <h3 class="modal-title">{{ currentMetric.title }}</h3>
          </div>
        </div>
        <button class="close-btn" @click="$emit('close')" aria-label="Close modal">&times;</button>
      </div>

      <div class="explainer-body">
        <!-- Quick Summary Box -->
        <div class="summary-highlight">
          <div class="highlight-label">In Simple Words:</div>
          <p class="highlight-text">{{ currentMetric.simpleExplanation }}</p>
        </div>

        <!-- Metric Target Range / Good vs Bad Indicator -->
        <div class="range-indicator">
          <div class="range-header">
            <span class="range-title">Target Benchmark</span>
            <span class="range-status" :class="currentMetric.benchmarkClass">
              {{ currentMetric.benchmarkText }}
            </span>
          </div>
          <div class="range-bar-track">
            <div class="range-bar-fill" :style="{ width: currentMetric.benchmarkProgress + '%' }"></div>
          </div>
          <div class="range-labels">
            <span>{{ currentMetric.rangeMin }}</span>
            <span>{{ currentMetric.rangeOptimal }}</span>
            <span>{{ currentMetric.rangeMax }}</span>
          </div>
        </div>

        <!-- Real World Example -->
        <div class="example-box">
          <div class="example-title">🎯 Real World Example</div>
          <p class="example-desc">{{ currentMetric.example }}</p>
        </div>

        <!-- Actionable Pro-Tip -->
        <div class="pro-tip">
          <span class="tip-icon">⚡</span>
          <div>
            <strong>Smart Investor Tip:</strong>
            <p>{{ currentMetric.actionableTip }}</p>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn btn-primary btn-block" @click="$emit('close')">
          Got it, back to Portfolio
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  metricKey: {
    type: String,
    default: 'sharpe'
  }
})

defineEmits(['close'])

const metricsData = {
  sharpe: {
    title: 'Sharpe Ratio',
    simpleExplanation: 'Measures how much return you earn per unit of risk taken compared to a safe asset (like government bonds). Higher is much better!',
    benchmarkText: 'Good ( > 1.0 ) | Excellent ( > 2.0 )',
    benchmarkClass: 'text-success',
    benchmarkProgress: 75,
    rangeMin: '< 0 (Poor)',
    rangeOptimal: '1.0 - 2.0 (Solid)',
    rangeMax: '3.0+ (Legendary)',
    example: 'If Portfolio A has 15% return with 10% volatility, and Portfolio B has 15% return with 30% volatility, Portfolio A has a much higher Sharpe ratio because it gave the same profit with much less stress.',
    actionableTip: 'Improve your Sharpe ratio by balancing high-growth tech stocks with defensive index funds or dividend assets.'
  },
  beta: {
    title: 'Portfolio Beta (β)',
    simpleExplanation: 'Measures how aggressively your portfolio moves relative to the broader market index (S&P 500 / SPY).',
    benchmarkText: 'Market baseline = 1.0',
    benchmarkClass: 'text-primary',
    benchmarkProgress: 55,
    rangeMin: '< 0.8 (Defensive)',
    rangeOptimal: '1.0 (Market Equal)',
    rangeMax: '> 1.3 (Aggressive)',
    example: 'If your Beta is 1.50 and the market goes up 2%, your portfolio is likely to jump ~3%. But if the market drops 2%, your portfolio is likely to drop ~3%.',
    actionableTip: 'Keep Beta close to 1.0 for steady long-term compounding, or lower it below 0.9 if you want lower downside swings.'
  },
  volatility: {
    title: 'Annualized Volatility',
    simpleExplanation: 'The standard deviation of daily price fluctuations scaled over a 252-day trading year. It reflects price turbulence.',
    benchmarkText: 'Moderate ( 12% - 20% )',
    benchmarkClass: 'text-warning',
    benchmarkProgress: 60,
    rangeMin: '5% (Calm / Bonds)',
    rangeOptimal: '15% (Balanced Stock)',
    rangeMax: '40%+ (Wild Swing)',
    example: 'High volatility means your portfolio net worth jumps up and down quickly like a roller coaster. Low volatility feels like a smooth train ride.',
    actionableTip: 'You can smooth volatility by spreading capital across non-correlated sectors (e.g. tech, health, finance, consumer).'
  },
  drawdown: {
    title: 'Maximum Drawdown (MDD)',
    simpleExplanation: 'The largest drop from a previous peak to a trough before a new peak is reached. It represents your worst-case historical drop.',
    benchmarkText: 'Controlled ( < 15% )',
    benchmarkClass: 'text-danger',
    benchmarkProgress: 40,
    rangeMin: '< 10% (Low Risk)',
    rangeOptimal: '15% - 25% (Standard Equity)',
    rangeMax: '> 40% (Extreme Loss)',
    example: 'If your portfolio went from $100k up to $150k and then plunged to $105k, your Maximum Drawdown was 30% ($45k drop from $150k peak).',
    actionableTip: 'Always keep an emergency cash buffer (10-20%) so you are never forced to sell assets at peak drawdown.'
  },
  var: {
    title: 'Value at Risk (VaR)',
    simpleExplanation: 'The maximum dollar loss expected over 1 trading day with a specific statistical confidence level (e.g. 95% or 99%).',
    benchmarkText: '95% & 99% 1-Day Horizon',
    benchmarkClass: 'text-danger',
    benchmarkProgress: 70,
    rangeMin: '95% (1 in 20 days)',
    rangeOptimal: 'Statistically Modeled',
    rangeMax: '99% (1 in 100 days)',
    example: 'A 95% 1-day VaR of $1,200 means there is a 95% chance you won\'t lose more than $1,200 tomorrow. Only in 1 out of 20 days might losses exceed that.',
    actionableTip: 'Institutional hedge funds use VaR to size their positions so a single bad day never wipes out the fund.'
  },
  health: {
    title: 'Portfolio Health Score',
    simpleExplanation: 'An aggregate 0-100 diagnostic evaluating your diversification, risk-adjusted returns (Sharpe), volatility moderation, and cash reserve.',
    benchmarkText: 'Target: 80+ / 100',
    benchmarkClass: 'text-success',
    benchmarkProgress: 85,
    rangeMin: '0 - 49 (High Risk)',
    rangeOptimal: '50 - 79 (Moderate)',
    rangeMax: '80 - 100 (Optimal)',
    example: 'A score of 85+ means your investments are well distributed across companies with healthy risk management principles.',
    actionableTip: 'Boost your health score by avoiding holding 100% in a single company and keeping at least 5% cash for opportunities.'
  }
}

const currentMetric = computed(() => {
  return metricsData[props.metricKey] || metricsData['sharpe']
})
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
  z-index: 1100;
  padding: 1rem;
}

.modal-content {
  width: 100%;
  max-width: 520px;
  padding: 1.75rem;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
  background: rgba(15, 23, 42, 0.95);
}

@media (max-width: 640px) {
  .modal-backdrop {
    align-items: flex-end;
    padding: 0;
  }

  .modal-content {
    border-radius: 24px 24px 0 0;
    max-height: 85vh;
    padding: 1.5rem 1.25rem 2rem 1.25rem;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.metric-icon-badge {
  font-size: 1.75rem;
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: rgba(99, 102, 241, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
}

.chip-sm {
  font-size: 0.65rem;
  font-weight: 700;
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
  padding: 0.25rem;
  line-height: 1;
}

.summary-highlight {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.12), rgba(139, 92, 246, 0.08));
  border: 1px solid rgba(99, 102, 241, 0.3);
  padding: 1rem;
  border-radius: 12px;
  margin-bottom: 1.25rem;
}

.highlight-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: #a5b4fc;
  text-transform: uppercase;
  margin-bottom: 0.25rem;
}

.highlight-text {
  font-size: 0.95rem;
  color: white;
  line-height: 1.45;
}

.range-indicator {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: 12px;
  margin-bottom: 1.25rem;
}

.range-header {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  margin-bottom: 0.5rem;
}

.range-title {
  color: var(--text-muted);
  font-weight: 600;
}

.range-status {
  font-weight: 700;
}

.range-bar-track {
  height: 8px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.4rem;
}

.range-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366f1, #10b981);
  border-radius: 4px;
}

.range-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.7rem;
  color: var(--text-dim);
}

.example-box {
  background: rgba(255, 255, 255, 0.025);
  border-left: 3px solid var(--primary);
  padding: 0.85rem 1rem;
  border-radius: 0 10px 10px 0;
  margin-bottom: 1.25rem;
}

.example-title {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 0.35rem;
}

.example-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  line-height: 1.45;
}

.pro-tip {
  display: flex;
  gap: 0.75rem;
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.25);
  padding: 0.85rem 1rem;
  border-radius: 12px;
  margin-bottom: 1.5rem;
  font-size: 0.825rem;
  color: #fef3c7;
}

.tip-icon {
  font-size: 1.2rem;
}

.modal-footer {
  margin-top: 0.5rem;
}

.btn-block {
  width: 100%;
  padding: 0.75rem;
}
</style>
