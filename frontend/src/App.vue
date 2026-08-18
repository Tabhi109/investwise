<template>
  <div class="app-wrapper">
    <!-- Toast Notifications -->
    <div class="toast-container" v-if="toasts.length > 0">
      <div 
        v-for="t in toasts" 
        :key="t.id" 
        class="toast"
        :class="{ 'badge-success': t.type === 'success', 'badge-danger': t.type === 'error' }"
      >
        {{ t.message }}
      </div>
    </div>

    <!-- Top Sticky Header -->
    <Navbar 
      :user="currentUser" 
      :activeTab="activeTab"
      @open-auth="showAuthModal = true" 
      @logout="handleLogout" 
      @update:activeTab="activeTab = $event"
    />

    <main class="main-content">
      <!-- Live Realtime Stock Feed Ticker Bar (Always Accessible) -->
      <TickerBar 
        :priceMap="priceMap" 
        @select-ticker="openTradeWithTicker" 
      />

      <!-- GUEST / LOGGED OUT BANNER (if user is not signed in) -->
      <div v-if="!currentUser" class="glass-card guest-hero-banner modal-animate">
        <div class="guest-banner-grid">
          <div class="guest-info">
            <div class="guest-badge">✨ NO-RISK SIMULATED INVESTING</div>
            <h1 class="guest-title">Learn. Trade. Analyze Risk. Master the Markets.</h1>
            <p class="guest-desc">
              InvestWise is not just a stock trading simulator — it is your personal quantitative investment coach. Receive <strong>$100,000</strong> in virtual funds, master risk analytics (Sharpe, Beta, VaR), test market shocks, and level up with interactive lessons.
            </p>
            <div class="guest-actions">
              <button @click="showAuthModal = true" class="btn btn-primary btn-lg" id="guest-get-started-btn">
                Claim $100,000 Demo Account
              </button>
              <button @click="activeTab = 'learn'" class="btn btn-secondary btn-lg">
                Explore Academy 🎓
              </button>
            </div>
          </div>
          <div class="guest-features-card">
            <div class="feat-item" @click="activeTab = 'markets'">
              <span class="feat-icon">⚡</span>
              <div>
                <strong>Live WebSocket Streaming</strong>
                <p>Sub-second tick feeds for mega-cap US equities & ETFs.</p>
              </div>
            </div>
            <div class="feat-item" @click="activeTab = 'risk'">
              <span class="feat-icon">🧠</span>
              <div>
                <strong>Quant Health & Stress Tester</strong>
                <p>Institutional covariance matrices, Sharpe Ratio, and shock simulations.</p>
              </div>
            </div>
            <div class="feat-item" @click="activeTab = 'learn'">
              <span class="feat-icon">🎓</span>
              <div>
                <strong>Smart Investor Academy</strong>
                <p>Interactive bite-sized courses that reward you with bonus simulated cash.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- MAIN APP TABS -->
      <div class="tab-content-area">
        <!-- TAB 1: PORTFOLIO DASHBOARD -->
        <section v-show="activeTab === 'portfolio'" class="tab-pane">
          <!-- Hero Portfolio Metrics -->
          <PortfolioOverview 
            :summary="portfolioSummary" 
            @open-trade="openTradeModal" 
            @explain-metric="openMetricExplainer"
          />

          <!-- Holdings Table & Mobile Cards -->
          <HoldingsTable 
            :holdings="portfolioSummary.holdings" 
            @open-trade="openTradeModal"
            @trade-item="handleTradeItem" 
          />

          <!-- Order Audit History -->
          <TransactionHistory 
            :transactions="transactions" 
          />
        </section>

        <!-- TAB 2: LIVE MARKETS & TRADING -->
        <section v-show="activeTab === 'markets'" class="tab-pane">
          <div class="markets-hub">
            <div class="hub-header">
              <div>
                <h2 class="section-heading">Live Market Feeds</h2>
                <p class="section-sub">Real-time quotes with instant 1-touch execution.</p>
              </div>
              <button @click="openTradeModal('BUY')" class="btn btn-primary btn-sm">
                + Create Custom Order
              </button>
            </div>

            <div class="markets-grid">
              <div 
                v-for="ticker in activeTickers" 
                :key="ticker"
                class="glass-card market-stock-card"
              >
                <div class="stock-card-top">
                  <div>
                    <span class="symbol-badge-lg">{{ ticker }}</span>
                    <span class="stock-company">{{ getCompanyName(ticker) }}</span>
                  </div>
                  <span :class="['stock-pct-chip', (priceMap[ticker]?.change || 0) >= 0 ? 'pct-up' : 'pct-down']">
                    {{ (priceMap[ticker]?.change || 0) >= 0 ? '+' : '' }}{{ (priceMap[ticker]?.change || 0).toFixed(2) }}%
                  </span>
                </div>

                <div class="stock-price-row">
                  <span class="stock-live-price font-heading">
                    ${{ (priceMap[ticker]?.price || 150).toFixed(2) }}
                  </span>
                </div>

                <div class="stock-meta-row">
                  <span>Category: {{ getTickerSector(ticker) }}</span>
                </div>

                <div class="stock-card-actions">
                  <button 
                    @click="openTradeWithTickerAndType(ticker, 'BUY')" 
                    class="btn btn-success btn-sm"
                  >
                    Buy {{ ticker }}
                  </button>
                  <button 
                    @click="openTradeWithTickerAndType(ticker, 'SELL')" 
                    class="btn btn-danger btn-sm"
                  >
                    Sell {{ ticker }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- TAB 3: QUANT COACH & RISK LAB -->
        <section v-show="activeTab === 'risk'" class="tab-pane">
          <RiskEngineDashboard 
            :metrics="riskMetrics" 
            :summary="portfolioSummary"
            :loading="riskLoading" 
            @refresh-risk="fetchRiskMetrics" 
            @explain-metric="openMetricExplainer"
          />
        </section>

        <!-- TAB 4: SMART INVESTOR ACADEMY -->
        <section v-show="activeTab === 'learn'" class="tab-pane">
          <LearnAcademy 
            @reward-user="handleQuizReward"
          />
        </section>
      </div>
    </main>

    <!-- Mobile Bottom Floating Navigation Bar -->
    <BottomNav 
      :activeTab="activeTab" 
      @update:activeTab="activeTab = $event"
    />

    <!-- MODALS -->
    <!-- Authentication Modal -->
    <AuthModal 
      v-if="showAuthModal" 
      :loading="authLoading" 
      :error="authError" 
      @close="showAuthModal = false" 
      @submit-auth="handleAuthSubmit" 
    />

    <!-- Trade Execution Modal / Bottom Sheet -->
    <TradeModal 
      v-if="showTradeModal" 
      :tickers="activeTickers" 
      :priceMap="priceMap" 
      :initialTicker="tradeModalTicker" 
      :initialType="tradeModalType" 
      :cashBalance="portfolioSummary.cash_balance" 
      :holdings="portfolioSummary.holdings"
      :loading="tradeLoading" 
      :error="tradeError" 
      @close="showTradeModal = false" 
      @execute-trade="handleTradeSubmit" 
    />

    <!-- Interactive Metric Explainer Modal -->
    <MetricExplainerModal 
      v-if="showExplainerModal"
      :metricKey="activeExplainerKey"
      @close="showExplainerModal = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import Navbar from './components/Navbar.vue'
import BottomNav from './components/BottomNav.vue'
import TickerBar from './components/TickerBar.vue'
import AuthModal from './components/AuthModal.vue'
import PortfolioOverview from './components/PortfolioOverview.vue'
import HoldingsTable from './components/HoldingsTable.vue'
import TradeModal from './components/TradeModal.vue'
import RiskEngineDashboard from './components/RiskEngineDashboard.vue'
import TransactionHistory from './components/TransactionHistory.vue'
import LearnAcademy from './components/LearnAcademy.vue'
import MetricExplainerModal from './components/MetricExplainerModal.vue'

// State
const activeTab = ref('portfolio') // 'portfolio' | 'markets' | 'risk' | 'learn'
const currentUser = ref(null)
const token = ref(localStorage.getItem('investwise_token') || '')

const activeTickers = ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'NVDA', 'SPY']
const priceMap = reactive({
  AAPL: { price: 182.5, change: 0.85, flash: null },
  MSFT: { price: 415.2, change: 1.20, flash: null },
  GOOGL: { price: 165.0, change: -0.45, flash: null },
  AMZN: { price: 185.3, change: 0.60, flash: null },
  NVDA: { price: 880.0, change: 2.30, flash: null },
  SPY: { price: 512.4, change: 0.40, flash: null }
})

const portfolioSummary = reactive({
  total_value: 100000.0,
  cash_balance: 100000.0,
  holdings_value: 0.0,
  unrealized_pnl: 0.0,
  realized_pnl: 0.0,
  holdings: []
})

const riskMetrics = ref(null)
const riskLoading = ref(false)
const transactions = ref([])

// Toast Notifications
const toasts = ref([])
function showToast(message, type = 'info') {
  const id = Date.now() + Math.random()
  toasts.value.push({ id, message, type })
  setTimeout(() => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }, 4000)
}

// Modals
const showAuthModal = ref(false)
const authLoading = ref(false)
const authError = ref('')

const showTradeModal = ref(false)
const tradeModalTicker = ref('AAPL')
const tradeModalType = ref('BUY')
const tradeLoading = ref(false)
const tradeError = ref('')

const showExplainerModal = ref(false)
const activeExplainerKey = ref('sharpe')

let ws = null

const companyNames = {
  AAPL: 'Apple Inc.',
  MSFT: 'Microsoft Corp.',
  GOOGL: 'Alphabet Inc.',
  AMZN: 'Amazon.com Inc.',
  NVDA: 'NVIDIA Corp.',
  SPY: 'S&P 500 ETF Trust'
}

function getCompanyName(ticker) {
  return companyNames[ticker] || 'Stock'
}

function getTickerSector(ticker) {
  if (ticker === 'SPY') return 'Broad Market ETF'
  if (ticker === 'NVDA') return 'Semiconductors'
  return 'Mega-Cap Technology'
}

onMounted(() => {
  if (token.value) {
    currentUser.value = { email: localStorage.getItem('investwise_user_email') || 'user@investwise.io' }
    fetchUserData()
  }
  connectWebSocket()
})

onUnmounted(() => {
  if (ws) ws.close()
})

// WebSocket Connection Management
function connectWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = window.location.port === '3000' 
    ? 'ws://localhost:8080/ws' 
    : `${protocol}//${window.location.host}/ws`

  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    activeTickers.forEach(ticker => {
      ws.send(JSON.stringify({ action: 'subscribe', ticker }))
    })
  }

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if (msg.type === 'ticker' && msg.ticker) {
        const prevPrice = priceMap[msg.ticker] ? priceMap[msg.ticker].price : msg.price
        const flashType = msg.price > prevPrice ? 'up' : (msg.price < prevPrice ? 'down' : null)
        
        priceMap[msg.ticker] = {
          price: msg.price,
          change: msg.change,
          flash: flashType
        }

        setTimeout(() => {
          if (priceMap[msg.ticker]) priceMap[msg.ticker].flash = null
        }, 800)
      }
    } catch (e) {
      // Ignore unparseable frames
    }
  }

  ws.onclose = () => {
    setTimeout(connectWebSocket, 3000)
  }
}

// User Data Fetching
async function fetchUserData() {
  if (!token.value) return
  await fetchPortfolioSummary()
  await fetchRiskMetrics()
  await fetchTransactions()
}

async function fetchPortfolioSummary() {
  try {
    const res = await fetch('/portfolio', {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    if (res.ok) {
      const data = await res.json()
      Object.assign(portfolioSummary, data)
    }
  } catch (e) {
    console.error('Failed to fetch portfolio summary', e)
  }
}

async function fetchRiskMetrics() {
  if (!token.value) return
  riskLoading.value = true
  try {
    const res = await fetch('/portfolio/risk', {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    if (res.ok) {
      riskMetrics.value = await res.json()
    }
  } catch (e) {
    console.error('Failed to fetch risk metrics', e)
  } finally {
    riskLoading.value = false
  }
}

async function fetchTransactions() {
  try {
    const res = await fetch('/portfolio/transactions', {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    if (res.ok) {
      transactions.value = await res.json()
    }
  } catch (e) {
    console.error('Failed to fetch transactions', e)
  }
}

// Authentication Handling
async function handleAuthSubmit(credentials) {
  authLoading.value = true
  authError.value = ''

  const endpoint = credentials.isRegister ? '/register' : '/login'

  try {
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials)
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || 'Authentication failed')
    }

    token.value = data.token
    localStorage.setItem('investwise_token', data.token)
    localStorage.setItem('investwise_user_email', data.email || credentials.email)
    currentUser.value = { email: data.email || credentials.email }

    showAuthModal.value = false
    showToast(`Welcome ${currentUser.value.email}! Your $100,000 demo capital is active.`, 'success')
    await fetchUserData()
  } catch (e) {
    authError.value = e.message
  } finally {
    authLoading.value = false
  }
}

function handleLogout() {
  token.value = ''
  currentUser.value = null
  localStorage.removeItem('investwise_token')
  localStorage.removeItem('investwise_user_email')
  portfolioSummary.total_value = 100000.0
  portfolioSummary.cash_balance = 100000.0
  portfolioSummary.holdings_value = 0.0
  portfolioSummary.unrealized_pnl = 0.0
  portfolioSummary.realized_pnl = 0.0
  portfolioSummary.holdings = []
  riskMetrics.value = null
  transactions.value = []
  showToast('Logged out successfully.')
}

// Trade Triggers
function openTradeModal(type = 'BUY') {
  tradeModalType.value = type
  showTradeModal.value = true
}

function openTradeWithTicker(ticker) {
  tradeModalTicker.value = ticker
  tradeModalType.value = 'BUY'
  showTradeModal.value = true
}

function openTradeWithTickerAndType(ticker, type) {
  tradeModalTicker.value = ticker
  tradeModalType.value = type
  showTradeModal.value = true
}

function handleTradeItem(item) {
  tradeModalTicker.value = item.ticker
  tradeModalType.value = item.type
  showTradeModal.value = true
}

async function handleTradeSubmit(tradePayload) {
  if (!token.value) {
    showAuthModal.value = true
    return
  }

  tradeLoading.value = true
  tradeError.value = ''

  try {
    const res = await fetch('/trade', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token.value}`
      },
      body: JSON.stringify(tradePayload)
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || 'Trade execution failed')
    }

    showTradeModal.value = false
    showToast(`Order executed: ${tradePayload.type} ${tradePayload.shares} shares of ${tradePayload.ticker}`, 'success')
    await fetchUserData()
  } catch (e) {
    tradeError.value = e.message
  } finally {
    tradeLoading.value = false
  }
}

// Metric Explainer Modal Trigger
function openMetricExplainer(metricKey) {
  activeExplainerKey.value = metricKey
  showExplainerModal.value = true
}

// Quiz Bonus Reward Handler
function handleQuizReward(amount) {
  portfolioSummary.cash_balance += amount
  portfolioSummary.total_value += amount
  showToast(`🎉 Milestone reached! +$${amount} added to your simulated portfolio!`, 'success')
}
</script>

<style scoped>
.app-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem 3rem 1rem;
}

/* Guest Banner */
.guest-hero-banner {
  padding: 2rem;
  margin-bottom: 1.5rem;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.14), rgba(139, 92, 246, 0.08));
  border-color: rgba(99, 102, 241, 0.35);
}

.guest-banner-grid {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 2rem;
  align-items: center;
}

@media (max-width: 900px) {
  .guest-banner-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
}

.guest-badge {
  display: inline-block;
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.4);
  font-size: 0.7rem;
  font-weight: 800;
  padding: 0.2rem 0.6rem;
  border-radius: 9999px;
  margin-bottom: 0.75rem;
}

.guest-title {
  font-size: 1.85rem;
  font-weight: 800;
  line-height: 1.25;
  margin-bottom: 0.75rem;
}

@media (max-width: 640px) {
  .guest-title {
    font-size: 1.5rem;
  }
}

.guest-desc {
  color: var(--text-muted);
  font-size: 0.95rem;
  line-height: 1.5;
  margin-bottom: 1.5rem;
}

.guest-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.guest-features-card {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.feat-item {
  display: flex;
  gap: 0.75rem;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 10px;
  transition: all 0.2s;
}

.feat-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.feat-icon {
  font-size: 1.35rem;
}

.feat-item strong {
  font-size: 0.875rem;
  color: white;
  display: block;
}

.feat-item p {
  font-size: 0.75rem;
  color: var(--text-dim);
}

.btn-lg {
  padding: 0.75rem 1.75rem;
  font-size: 0.95rem;
  border-radius: 12px;
}

/* Tab Content Areas */
.tab-content-area {
  margin-top: 0.5rem;
}

/* Markets Tab Styles */
.markets-hub {
  margin-bottom: 2rem;
}

.hub-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.25rem;
}

.section-heading {
  font-size: 1.35rem;
  font-weight: 800;
}

.section-sub {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.markets-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1rem;
}

.market-stock-card {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.stock-card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.75rem;
}

.symbol-badge-lg {
  font-size: 1.1rem;
  font-weight: 800;
  color: #a5b4fc;
  display: block;
}

.stock-company {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.stock-pct-chip {
  font-size: 0.8rem;
  font-weight: 700;
  padding: 0.2rem 0.5rem;
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

.stock-price-row {
  margin-bottom: 0.5rem;
}

.stock-live-price {
  font-size: 1.65rem;
  font-weight: 800;
}

.stock-meta-row {
  font-size: 0.75rem;
  color: var(--text-dim);
  margin-bottom: 1rem;
}

.stock-card-actions {
  display: flex;
  gap: 0.5rem;
}

.stock-card-actions button {
  flex: 1;
}
</style>
