<template>
  <div class="app-wrapper">
    <Navbar 
      :user="currentUser" 
      @open-auth="showAuthModal = true" 
      @logout="handleLogout" 
    />

    <main class="main-content">
      <!-- Live Realtime Stock Feed Ticker Bar -->
      <TickerBar 
        :priceMap="priceMap" 
        @select-ticker="openTradeWithTicker" 
      />

      <template v-if="currentUser">
        <!-- Hero Portfolio Metrics -->
        <PortfolioOverview 
          :summary="portfolioSummary" 
          @open-trade="openTradeModal" 
        />

        <!-- Holdings Table -->
        <HoldingsTable 
          :holdings="portfolioSummary.holdings" 
          @open-trade="openTradeModal"
          @trade-item="handleTradeItem" 
        />

        <!-- Quantitative Risk Engine Section -->
        <RiskEngineDashboard 
          :metrics="riskMetrics" 
          :loading="riskLoading" 
          @refresh-risk="fetchRiskMetrics" 
        />

        <!-- Order Audit History -->
        <TransactionHistory 
          :transactions="transactions" 
        />
      </template>

      <div v-else class="glass-card guest-banner">
        <div class="banner-content">
          <h2 class="banner-title">Real-time Stock Trading & Quantitative Portfolio Analytics</h2>
          <p class="banner-desc">
            Sign in or register to receive $100,000 in simulated cash, place real-time trades, and compute quantitative metrics like Sharpe Ratio, Beta, Max Drawdown, and Value at Risk (VaR).
          </p>
          <button @click="showAuthModal = true" class="btn btn-primary btn-lg" id="guest-get-started-btn">
            Get Started with $100,000 Demo Account
          </button>
        </div>
      </div>
    </main>

    <!-- Modals -->
    <AuthModal 
      v-if="showAuthModal" 
      :loading="authLoading" 
      :error="authError" 
      @close="showAuthModal = false" 
      @submit-auth="handleAuthSubmit" 
    />

    <TradeModal 
      v-if="showTradeModal" 
      :tickers="activeTickers" 
      :priceMap="priceMap" 
      :initialTicker="tradeModalTicker" 
      :initialType="tradeModalType" 
      :cashBalance="portfolioSummary.cash_balance" 
      :loading="tradeLoading" 
      :error="tradeError" 
      @close="showTradeModal = false" 
      @execute-trade="handleTradeSubmit" 
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import Navbar from './components/Navbar.vue'
import TickerBar from './components/TickerBar.vue'
import AuthModal from './components/AuthModal.vue'
import PortfolioOverview from './components/PortfolioOverview.vue'
import HoldingsTable from './components/HoldingsTable.vue'
import TradeModal from './components/TradeModal.vue'
import RiskEngineDashboard from './components/RiskEngineDashboard.vue'
import TransactionHistory from './components/TransactionHistory.vue'

// State
const currentUser = ref(null)
const token = ref(localStorage.getItem('investwise_token') || '')

const activeTickers = ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'NVDA', 'SPY']
const priceMap = reactive({
  AAPL: { price: 180.0, change: 0.0, flash: null },
  MSFT: { price: 400.0, change: 0.0, flash: null },
  GOOGL: { price: 150.0, change: 0.0, flash: null },
  AMZN: { price: 175.0, change: 0.0, flash: null },
  NVDA: { price: 800.0, change: 0.0, flash: null },
  SPY: { price: 500.0, change: 0.0, flash: null }
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

// Modals
const showAuthModal = ref(false)
const authLoading = ref(false)
const authError = ref('')

const showTradeModal = ref(false)
const tradeModalTicker = ref('AAPL')
const tradeModalType = ref('BUY')
const tradeLoading = ref(false)
const tradeError = ref('')

let ws = null

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
    // Subscribe to all tickers
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

        // Reset flash state
        setTimeout(() => {
          if (priceMap[msg.ticker]) priceMap[msg.ticker].flash = null
        }, 800)
      }
    } catch (e) {
      // Ignore unparseable frames
    }
  }

  ws.onclose = () => {
    setTimeout(connectWebSocket, 3000) // Reconnect after 3s
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

// Authentication Handlers
async function handleAuthSubmit(payload) {
  authLoading.value = true
  authError.value = ''

  const endpoint = payload.isRegister ? '/register' : '/login'
  const body = payload.isRegister 
    ? { username: payload.username, email: payload.email, password: payload.password }
    : { email: payload.email, password: payload.password }

  try {
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || 'Authentication failed')
    }

    if (payload.isRegister) {
      // Auto login after register
      await handleAuthSubmit({ isRegister: false, email: payload.email, password: payload.password })
      return
    }

    token.value = data.token
    localStorage.setItem('investwise_token', data.token)
    localStorage.setItem('investwise_user_email', payload.email)
    currentUser.value = { email: payload.email }

    showAuthModal.value = false
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
  riskMetrics.value = null
}

// Trade Handlers
function openTradeModal(type = 'BUY') {
  tradeModalType.value = type
  tradeModalTicker.value = 'AAPL'
  showTradeModal.value = true
}

function openTradeWithTicker(ticker) {
  tradeModalTicker.value = ticker
  tradeModalType.value = 'BUY'
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
    await fetchUserData()
  } catch (e) {
    tradeError.value = e.message
  } finally {
    tradeLoading.value = false
  }
}
</script>

<style scoped>
.app-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem 3rem 1rem;
}

.guest-banner {
  padding: 4rem 2rem;
  text-align: center;
  margin-top: 2rem;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.12), rgba(139, 92, 246, 0.08));
  border-color: rgba(99, 102, 241, 0.3);
}

.banner-content {
  max-width: 700px;
  margin: 0 auto;
}

.banner-title {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 1rem;
  line-height: 1.25;
}

.banner-desc {
  color: var(--text-muted);
  font-size: 1.05rem;
  margin-bottom: 2rem;
}

.btn-lg {
  padding: 0.85rem 2rem;
  font-size: 1rem;
  border-radius: 12px;
}
</style>
