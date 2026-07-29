<template>
  <div class="glass-card risk-card">
    <div class="risk-header">
      <div class="title-wrap">
        <div class="chip">QUANT ENGINE</div>
        <h3 class="card-title">Portfolio Quantitative Analytics</h3>
      </div>
      <button @click="$emit('refresh-risk')" class="btn btn-secondary btn-sm" id="refresh-risk-btn">
        Refresh Risk Analytics
      </button>
    </div>

    <div v-if="loading" class="loading-state">
      Computing covariance matrices & statistical models...
    </div>

    <div v-else-if="!metrics" class="empty-state">
      Log in and hold stock positions to generate quantitative risk analytics.
    </div>

    <div v-else class="metrics-content">
      <!-- Top Summary Grid -->
      <div class="risk-stats-grid">
        <div class="stat-box">
          <div class="stat-title">Sharpe Ratio</div>
          <div class="stat-num" :class="metrics.sharpe_ratio >= 1 ? 'text-success' : (metrics.sharpe_ratio >= 0 ? 'text-warning' : 'text-danger')">
            {{ formatNum(metrics.sharpe_ratio, 2) }}
          </div>
          <div class="stat-sub">Risk-adjusted return vs 4% R_f</div>
        </div>

        <div class="stat-box">
          <div class="stat-title">Portfolio Beta (β)</div>
          <div class="stat-num text-primary">
            {{ formatNum(metrics.beta, 2) }}
          </div>
          <div class="stat-sub">Market sensitivity vs SPY</div>
        </div>

        <div class="stat-box">
          <div class="stat-title">Annualized Volatility</div>
          <div class="stat-num">
            {{ (metrics.volatility * 100).toFixed(1) }}%
          </div>
          <div class="stat-sub">Standard deviation (252 days)</div>
        </div>

        <div class="stat-box">
          <div class="stat-title">Max Drawdown (MDD)</div>
          <div class="stat-num text-danger">
            {{ (metrics.max_drawdown * 100).toFixed(1) }}%
          </div>
          <div class="stat-sub">Worst historical peak-to-trough</div>
        </div>
      </div>

      <!-- Value at Risk (VaR) Panel -->
      <div class="var-panel">
        <h4 class="section-title">Value at Risk (1-Day Horizon)</h4>
        <div class="var-grid">
          <div class="var-card">
            <span class="var-type">Parametric (95% Conf)</span>
            <span class="var-val">${{ formatNum(metrics.parametric_var_95, 2) }}</span>
          </div>
          <div class="var-card">
            <span class="var-type">Parametric (99% Conf)</span>
            <span class="var-val text-danger">${{ formatNum(metrics.parametric_var_99, 2) }}</span>
          </div>
          <div class="var-card">
            <span class="var-type">Historical (95% Conf)</span>
            <span class="var-val">${{ formatNum(metrics.historical_var_95, 2) }}</span>
          </div>
          <div class="var-card">
            <span class="var-type">Historical (99% Conf)</span>
            <span class="var-val text-danger">${{ formatNum(metrics.historical_var_99, 2) }}</span>
          </div>
        </div>
      </div>

      <!-- Marginal Risk Contribution Bars -->
      <div v-if="metrics.risk_contributions && Object.keys(metrics.risk_contributions).length > 0" class="risk-contrib-panel">
        <h4 class="section-title">Marginal Risk Contribution per Asset</h4>
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
</template>

<script setup>
defineProps({
  metrics: Object,
  loading: Boolean
})

defineEmits(['refresh-risk'])

function formatNum(val, decimals = 2) {
  if (val === undefined || val === null) return '0.00'
  return val.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
}
</script>

<style scoped>
.risk-card {
  margin-bottom: 1.75rem;
  padding: 1.5rem;
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
  gap: 0.75rem;
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

.loading-state, .empty-state {
  padding: 2.5rem;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.risk-stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.stat-box {
  background: rgba(255, 255, 255, 0.025);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: 12px;
}

.stat-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}

.stat-num {
  font-family: var(--font-heading);
  font-size: 1.6rem;
  font-weight: 800;
  margin: 0.2rem 0;
}

.stat-sub {
  font-size: 0.7rem;
  color: var(--text-dim);
}

.text-success { color: var(--success); }
.text-warning { color: var(--warning); }
.text-danger { color: var(--danger); }
.text-primary { color: var(--primary); }

.section-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 0.75rem;
}

.var-panel {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1.25rem;
}

.var-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 0.75rem;
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
  margin-top: 0.15rem;
}

.contrib-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.contrib-header {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  margin-bottom: 0.25rem;
}

.contrib-ticker {
  font-weight: 700;
  color: #a5b4fc;
}

.contrib-vol {
  color: var(--text-muted);
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
