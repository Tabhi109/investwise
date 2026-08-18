<template>
  <div class="risk-section">
    <!-- Main Quant Analysis Card -->
    <div class="glass-card risk-card">
      <div class="risk-header">
        <div class="title-wrap">
          <div class="chip">QUANT ENGINE</div>
          <h3 class="card-title">Portfolio Quant Coach</h3>
        </div>
        <div class="header-actions">
          <button @click="$emit('refresh-risk')" class="btn btn-secondary btn-sm" id="refresh-risk-btn">
            🔄 Refresh Models
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <span class="spinner"></span>
        <p>Computing covariance matrices & statistical models...</p>
      </div>

      <div v-else-if="!metrics" class="empty-state">
        <div class="empty-icon">🛡️</div>
        <h4>No Positions for Risk Analysis</h4>
        <p>Purchase stocks to generate real-time quantitative risk analytics, Sharpe Ratios, Betas, and Value at Risk (VaR) models.</p>
      </div>

      <div v-else class="metrics-content">
        <!-- AI Quant Coach Insight Banner -->
        <div class="coach-banner">
          <div class="coach-avatar">🧠</div>
          <div class="coach-body">
            <div class="coach-title">Smart Quant Diagnostic</div>
            <p class="coach-text">{{ coachDiagnosticText }}</p>
          </div>
        </div>

        <!-- 4 Key Quant Metrics Grid (with Explainer Tap Triggers) -->
        <div class="risk-stats-grid">
          <!-- Sharpe Ratio -->
          <div class="stat-box clickable-stat" @click="$emit('explain-metric', 'sharpe')">
            <div class="stat-top">
              <span class="stat-title">Sharpe Ratio</span>
              <span class="info-badge">Explain ℹ️</span>
            </div>
            <div class="stat-num" :class="metrics.sharpe_ratio >= 1 ? 'text-success' : (metrics.sharpe_ratio >= 0 ? 'text-warning' : 'text-danger')">
              {{ formatNum(metrics.sharpe_ratio, 2) }}
            </div>
            <div class="stat-sub">Risk-adjusted return vs 4% R_f</div>
          </div>

          <!-- Beta -->
          <div class="stat-box clickable-stat" @click="$emit('explain-metric', 'beta')">
            <div class="stat-top">
              <span class="stat-title">Portfolio Beta (β)</span>
              <span class="info-badge">Explain ℹ️</span>
            </div>
            <div class="stat-num text-primary">
              {{ formatNum(metrics.beta, 2) }}
            </div>
            <div class="stat-sub">Sensitivity vs S&P 500 Index</div>
          </div>

          <!-- Volatility -->
          <div class="stat-box clickable-stat" @click="$emit('explain-metric', 'volatility')">
            <div class="stat-top">
              <span class="stat-title">Volatility</span>
              <span class="info-badge">Explain ℹ️</span>
            </div>
            <div class="stat-num text-warning">
              {{ (metrics.volatility * 100).toFixed(1) }}%
            </div>
            <div class="stat-sub">Annualized standard deviation</div>
          </div>

          <!-- Max Drawdown -->
          <div class="stat-box clickable-stat" @click="$emit('explain-metric', 'drawdown')">
            <div class="stat-top">
              <span class="stat-title">Max Drawdown</span>
              <span class="info-badge">Explain ℹ️</span>
            </div>
            <div class="stat-num text-danger">
              {{ (metrics.max_drawdown * 100).toFixed(1) }}%
            </div>
            <div class="stat-sub">Worst peak-to-trough drop</div>
          </div>
        </div>

        <!-- Value at Risk (VaR) Panel -->
        <div class="var-panel">
          <div class="var-header">
            <h4 class="section-title">Value at Risk (1-Day Horizon)</h4>
            <button class="btn btn-secondary btn-xs" @click="$emit('explain-metric', 'var')">
              What is VaR? ℹ️
            </button>
          </div>
          <div class="var-grid">
            <div class="var-card">
              <span class="var-type">Parametric (95% Conf)</span>
              <span class="var-val">${{ formatNum(metrics.parametric_var_95, 2) }}</span>
              <span class="var-desc">1 in 20 days worst case</span>
            </div>
            <div class="var-card">
              <span class="var-type">Parametric (99% Conf)</span>
              <span class="var-val text-danger">${{ formatNum(metrics.parametric_var_99, 2) }}</span>
              <span class="var-desc">1 in 100 days worst case</span>
            </div>
            <div class="var-card">
              <span class="var-type">Historical (95% Conf)</span>
              <span class="var-val">${{ formatNum(metrics.historical_var_95, 2) }}</span>
              <span class="var-desc">Empirical price distribution</span>
            </div>
            <div class="var-card">
              <span class="var-type">Historical (99% Conf)</span>
              <span class="var-val text-danger">${{ formatNum(metrics.historical_var_99, 2) }}</span>
              <span class="var-desc">Empirical tail risk</span>
            </div>
          </div>
        </div>

        <!-- Marginal Risk Contribution Bars -->
        <div v-if="metrics.risk_contributions && Object.keys(metrics.risk_contributions).length > 0" class="risk-contrib-panel">
          <h4 class="section-title">Risk Contribution per Asset (Volatility Driver)</h4>
          <div class="contrib-list">
            <div v-for="(volContrib, ticker) in metrics.risk_contributions" :key="ticker" class="contrib-item">
              <div class="contrib-header">
                <span class="contrib-ticker">{{ ticker }}</span>
                <span class="contrib-vol">+{{ (volContrib * 100).toFixed(2) }}% Volatility</span>
              </div>
              <div class="bar-bg">
                <div class="bar-fill" :style="{ width: Math.min(100, Math.max(5, (volContrib / (metrics.volatility || 1)) * 100)) + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Embedded Interactive Stress Tester -->
    <StressTester 
      :summary="summary" 
      :beta="metrics ? metrics.beta : 1.0" 
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import StressTester from './StressTester.vue'

const props = defineProps({
  metrics: Object,
  summary: Object,
  loading: Boolean
})

defineEmits(['refresh-risk', 'explain-metric'])

const coachDiagnosticText = computed(() => {
  if (!props.metrics) return 'Add positions to receive AI-driven portfolio health coaching.'

  const m = props.metrics
  if (m.sharpe_ratio >= 1.5) {
    return `Excellent risk-adjusted returns! Your Sharpe Ratio (${m.sharpe_ratio.toFixed(2)}) indicates your profits heavily outweigh the risk taken.`
  } else if (m.sharpe_ratio >= 0.8) {
    return `Solid portfolio foundation. Your Beta of ${m.beta.toFixed(2)} means you track the broader market closely. Consider adding index funds to smooth drawdowns.`
  } else if (m.beta > 1.3) {
    return `High Beta alert (${m.beta.toFixed(2)}). Your portfolio swings significantly more than the broader S&P 500. Ensure you maintain adequate cash reserves.`
  }
  return `Your portfolio is active. Maintain disciplined position sizing and diversify across multiple sectors to keep 1-day VaR within comfortable limits.`
})

function formatNum(val, decimals = 2) {
  if (val === undefined || val === null) return '0.00'
  return val.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
}
</script>

<style scoped>
.risk-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.risk-card {
  padding: 1.5rem;
}

@media (max-width: 640px) {
  .risk-card {
    padding: 1.25rem;
  }
}

.risk-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
}

.title-wrap {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.chip {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(139, 92, 246, 0.2));
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.4);
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

.btn-sm {
  padding: 0.4rem 0.85rem;
  font-size: 0.8rem;
}

.btn-xs {
  padding: 0.25rem 0.5rem;
  font-size: 0.7rem;
  border-radius: 6px;
}

/* Coach Banner */
.coach-banner {
  display: flex;
  gap: 0.85rem;
  align-items: flex-start;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.12), rgba(6, 182, 212, 0.08));
  border: 1px solid rgba(99, 102, 241, 0.3);
  padding: 1rem;
  border-radius: 12px;
  margin-bottom: 1.25rem;
}

.coach-avatar {
  font-size: 1.5rem;
  background: rgba(99, 102, 241, 0.2);
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.coach-title {
  font-size: 0.8rem;
  font-weight: 700;
  color: #a5b4fc;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: 0.2rem;
}

.coach-text {
  font-size: 0.875rem;
  color: #f1f5f9;
  line-height: 1.45;
}

/* Loading & Empty States */
.loading-state, .empty-state {
  padding: 2.5rem 1rem;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.empty-state h4 {
  font-size: 1.1rem;
  margin-bottom: 0.4rem;
  color: var(--text-main);
}

.empty-state p {
  max-width: 450px;
  margin: 0 auto;
  line-height: 1.45;
}

/* Stats Grid */
.risk-stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 0.85rem;
  margin-bottom: 1.25rem;
}

.stat-box {
  background: rgba(255, 255, 255, 0.025);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: 12px;
  transition: all 0.2s;
}

.clickable-stat {
  cursor: pointer;
}

.clickable-stat:hover {
  background: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.4);
  transform: translateY(-2px);
}

.stat-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-title {
  font-size: 0.725rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}

.info-badge {
  font-size: 0.65rem;
  color: #a5b4fc;
}

.stat-num {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  font-weight: 800;
  margin: 0.2rem 0;
}

.stat-sub {
  font-size: 0.68rem;
  color: var(--text-dim);
}

.text-success { color: var(--success); }
.text-warning { color: var(--warning); }
.text-danger { color: var(--danger); }
.text-primary { color: var(--primary); }

.section-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-main);
}

.var-panel {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1.25rem;
}

.var-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.var-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 0.65rem;
}

.var-card {
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.03);
  padding: 0.65rem 0.85rem;
  border-radius: 8px;
}

.var-type {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.var-val {
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 1.1rem;
  margin: 0.15rem 0;
}

.var-desc {
  font-size: 0.65rem;
  color: var(--text-dim);
}

.risk-contrib-panel {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 1rem;
}

.contrib-list {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
  margin-top: 0.75rem;
}

.contrib-header {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  margin-bottom: 0.2rem;
}

.contrib-ticker {
  font-weight: 700;
  color: #a5b4fc;
}

.contrib-vol {
  color: var(--text-muted);
  font-size: 0.75rem;
}

.bar-bg {
  height: 6px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), var(--accent));
  border-radius: 3px;
}
</style>
