<template>
  <div class="academy-container">
    <!-- Header Banner -->
    <div class="glass-card academy-hero">
      <div class="hero-content">
        <div class="hero-badge">
          <span class="pulse-spark">✨</span> SMART INVESTOR ACADEMY
        </div>
        <h2 class="hero-title">Learn by Doing: Master Modern Investing</h2>
        <p class="hero-desc">
          Bite-sized visual lessons, interactive quizzes, and risk simulators. Complete lessons to earn bonus simulated cash and level up your investment mastery!
        </p>

        <div class="progress-stats">
          <div class="stat-pill">
            <span class="pill-label">Completed Lessons</span>
            <span class="pill-val">{{ completedTracksCount }} / {{ tracks.length }}</span>
          </div>
          <div class="stat-pill">
            <span class="pill-label">Investor Level</span>
            <span class="pill-val text-accent">{{ currentLevelTitle }}</span>
          </div>
          <div class="stat-pill">
            <span class="pill-label">Earned Bonus Cash</span>
            <span class="pill-val text-success">+${{ earnedBonusTotal }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Active Lesson / Quiz View Modal or Inline Card -->
    <div v-if="activeLesson" class="glass-card active-lesson-card modal-animate">
      <div class="lesson-header">
        <div class="lesson-meta">
          <span class="badge badge-primary">{{ activeLesson.trackCategory }}</span>
          <h3 class="lesson-title">{{ activeLesson.title }}</h3>
        </div>
        <button class="btn btn-secondary btn-sm" @click="activeLesson = null">Close Lesson</button>
      </div>

      <div class="lesson-content-body">
        <div class="lesson-illustration">
          <span class="illus-icon">{{ activeLesson.icon }}</span>
        </div>

        <div class="lesson-text-block">
          <h4 class="concept-heading">{{ activeLesson.keyConceptTitle }}</h4>
          <p class="concept-p">{{ activeLesson.content }}</p>

          <!-- Interactive Micro-Tool: e.g. Compound Calculator if track is compounding -->
          <div v-if="activeLesson.id === 'compounding'" class="mini-calc-box">
            <div class="calc-header">
              <span>📈 The Power of 10% Annual Compounding over 20 Years</span>
            </div>
            <div class="calc-slider-wrap">
              <label>Monthly Contribution: <strong>${{ monthlyInvestment }}</strong></label>
              <input type="range" min="100" max="2000" step="50" v-model.number="monthlyInvestment" class="calc-slider" />
            </div>
            <div class="calc-result">
              <span>Future Net Worth: <strong>${{ calculateCompounding(monthlyInvestment).toLocaleString() }}</strong></span>
              <span class="text-success">(Total Invested: ${{ (monthlyInvestment * 12 * 20).toLocaleString() }} | Profit: ${{ (calculateCompounding(monthlyInvestment) - monthlyInvestment * 12 * 20).toLocaleString() }})</span>
            </div>
          </div>

          <div class="key-takeaways">
            <span class="takeaway-title">🔑 Core Rule to Remember:</span>
            <p>{{ activeLesson.goldenRule }}</p>
          </div>
        </div>

        <!-- Quiz Section -->
        <div class="quiz-section">
          <h4 class="quiz-title">🎯 Quick Knowledge Check (Reward: +${{ activeLesson.rewardCash }} Simulated Cash)</h4>
          <p class="quiz-question">{{ activeLesson.quiz.question }}</p>

          <div class="quiz-options">
            <button 
              v-for="(option, idx) in activeLesson.quiz.options" 
              :key="idx"
              class="quiz-opt-btn"
              :class="{
                'opt-selected': selectedAnswer === idx,
                'opt-correct': quizSubmitted && idx === activeLesson.quiz.correctIndex,
                'opt-wrong': quizSubmitted && selectedAnswer === idx && idx !== activeLesson.quiz.correctIndex
              }"
              :disabled="quizSubmitted"
              @click="selectedAnswer = idx"
            >
              <span class="opt-letter">{{ ['A', 'B', 'C', 'D'][idx] }}</span>
              <span class="opt-text">{{ option }}</span>
            </button>
          </div>

          <div v-if="quizSubmitted" class="quiz-feedback" :class="isAnswerCorrect ? 'feedback-success' : 'feedback-fail'">
            <div class="feedback-title">
              {{ isAnswerCorrect ? '🎉 Correct! Well Done!' : '❌ Not Quite!' }}
            </div>
            <p class="feedback-desc">{{ activeLesson.quiz.explanation }}</p>
            <div v-if="isAnswerCorrect && !activeLesson.claimed" class="reward-notice">
              <span>💰 +${{ activeLesson.rewardCash }} has been credited to your Demo account!</span>
            </div>
          </div>

          <div class="quiz-actions">
            <button 
              v-if="!quizSubmitted" 
              class="btn btn-primary btn-block"
              :disabled="selectedAnswer === null"
              @click="submitQuiz"
            >
              Submit Answer
            </button>
            <button 
              v-else 
              class="btn btn-secondary btn-block"
              @click="finishLesson"
            >
              Continue to Next Lesson
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Lesson Tracks Grid -->
    <div class="tracks-grid">
      <div 
        v-for="track in tracks" 
        :key="track.id"
        class="glass-card track-card"
        :class="{ 'track-completed': completedTracks[track.id] }"
        @click="openLesson(track)"
      >
        <div class="track-top">
          <div class="track-icon-wrap">{{ track.icon }}</div>
          <span v-if="completedTracks[track.id]" class="badge badge-success">✓ Completed</span>
          <span v-else class="badge badge-primary">+${{ track.rewardCash }} Reward</span>
        </div>

        <h3 class="track-title">{{ track.title }}</h3>
        <p class="track-summary">{{ track.summary }}</p>

        <div class="track-footer">
          <span class="track-meta">{{ track.duration }} read • {{ track.level }}</span>
          <button class="btn btn-sm" :class="completedTracks[track.id] ? 'btn-secondary' : 'btn-primary'">
            {{ completedTracks[track.id] ? 'Review' : 'Start Lesson' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Achievement Badges Showcase -->
    <div class="glass-card badges-card">
      <h3 class="card-title">🏆 Your Investor Badges</h3>
      <div class="badges-grid">
        <div 
          v-for="badge in badges" 
          :key="badge.id"
          class="badge-item"
          :class="{ 'badge-unlocked': badge.unlocked }"
        >
          <div class="badge-icon-box">{{ badge.icon }}</div>
          <span class="badge-name">{{ badge.name }}</span>
          <span class="badge-crit">{{ badge.criteria }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['reward-user'])

const monthlyInvestment = ref(500)

const tracks = [
  {
    id: 'compounding',
    trackCategory: 'FOUNDATIONS',
    icon: '🌱',
    title: 'The Eighth Wonder: Compound Growth',
    summary: 'How consistent investing turns modest regular contributions into life-changing exponential portfolios.',
    duration: '2 min',
    level: 'Beginner',
    rewardCash: 1000,
    keyConceptTitle: 'Exponential Growth vs Linear Savings',
    content: 'When you earn returns on both your initial investment and on accumulated interest from previous years, your money snowballs. Time in the market is vastly more powerful than timing the market.',
    goldenRule: 'Starting 5 years earlier can double your eventual retirement portfolio even if you invest less total capital.',
    quiz: {
      question: 'Why is compound interest considered more powerful than simple interest?',
      options: [
        'It requires no tax payments ever.',
        'You earn returns on both your principal and all previous accumulated gains.',
        'It only works if you buy high-risk penny stocks.',
        'It doubles your money every 30 days automatically.'
      ],
      correctIndex: 1,
      explanation: 'Compounding works because your returns generate their own returns, creating an accelerating upward growth curve over time.'
    }
  },
  {
    id: 'diversification',
    trackCategory: 'RISK MANAGEMENT',
    icon: '🛡️',
    title: 'Diversification & The Free Lunch',
    summary: 'Why Nobel Prize economists call diversification the only free lunch in finance, and how to eliminate unsystematic risk.',
    duration: '3 min',
    level: 'Beginner',
    rewardCash: 1000,
    keyConceptTitle: 'Unsystematic vs Systematic Risk',
    content: 'Holding only 1 or 2 tech stocks exposes you to catastrophic single-company collapse. By distributing capital across sectors (Index ETFs like SPY, tech, healthcare, cash), you reduce volatility without sacrificing expected market returns.',
    goldenRule: 'Never let a single stock position exceed 15-20% of your total liquid net worth.',
    quiz: {
      question: 'What is the primary benefit of holding a broad index fund like SPY along with individual stocks?',
      options: [
        'It eliminates all risk of losing money on any single day.',
        'It eliminates single-company specific risk and provides broad economic participation.',
        'It guarantees a fixed 50% return every year.',
        'It prevents you from paying broker fees.'
      ],
      correctIndex: 1,
      explanation: 'Holding broad index ETFs protects your portfolio from individual corporate failures while riding long-term economic expansion.'
    }
  },
  {
    id: 'beta_sharpe',
    trackCategory: 'QUANTITATIVE MASTERY',
    icon: '⚡',
    title: 'Mastering Sharpe Ratio & Beta',
    summary: 'Learn the exact math institutional hedge funds use to evaluate whether a trader is skilled or just taking reckless risks.',
    duration: '3 min',
    level: 'Intermediate',
    rewardCash: 1000,
    keyConceptTitle: 'Risk-Adjusted Returns',
    content: 'A high return is meaningless if it required extreme volatility and high risk of ruin. The Sharpe Ratio compares your excess return over the risk-free rate (e.g. 4% Treasuries) against your volatility. A Sharpe > 1.0 indicates good skill; > 2.0 is elite.',
    goldenRule: 'Optimize for Sharpe Ratio (risk-adjusted return) rather than chasing raw high-volatility spikes.',
    quiz: {
      question: 'If Portfolio A and Portfolio B both made 20% this year, but Portfolio A had half the volatility of B, which has the higher Sharpe Ratio?',
      options: [
        'Portfolio B has the higher Sharpe Ratio.',
        'Portfolio A has the higher Sharpe Ratio.',
        'Both have the identical Sharpe Ratio.',
        'Neither, Sharpe Ratio only applies to bonds.'
      ],
      correctIndex: 1,
      explanation: 'Portfolio A achieved the same return with half the volatility, giving it a much higher risk-adjusted efficiency score!'
    }
  },
  {
    id: 'var_drawdown',
    trackCategory: 'PORTFOLIO DEFENSE',
    icon: '📊',
    title: 'Value at Risk (VaR) & Managing Drawdowns',
    summary: 'How to calculate your worst-case daily loss and protect your psychological capital during market crashes.',
    duration: '3 min',
    level: 'Advanced',
    rewardCash: 1000,
    keyConceptTitle: 'Downside Protection & Position Sizing',
    content: 'Value at Risk (VaR 95%) informs you of the maximum dollar loss expected on 19 out of 20 trading days. Keeping an adequate cash reserve ensures you never panic-sell quality companies at bottom prices.',
    goldenRule: 'Loss aversion causes most retail investors to sell at the exact bottom. Having a clear risk plan prevents emotional trading.',
    quiz: {
      question: 'What does a 1-day 95% Value at Risk (VaR) of $1,500 mean?',
      options: [
        'You are guaranteed to lose $1,500 every single day.',
        'There is a 95% probability that your 1-day loss will not exceed $1,500.',
        'You can only make $1,500 maximum profit tomorrow.',
        'Your broker will liquidate your account if losses hit $1,500.'
      ],
      correctIndex: 1,
      explanation: 'VaR 95% models that on 95% of trading days, your daily drawdown will stay within the $1,500 threshold.'
    }
  }
]

// State & Persistence
const completedTracks = ref(JSON.parse(localStorage.getItem('investwise_completed_lessons') || '{}'))
const activeLesson = ref(null)
const selectedAnswer = ref(null)
const quizSubmitted = ref(false)

const completedTracksCount = computed(() => Object.keys(completedTracks.value).length)

const earnedBonusTotal = computed(() => {
  let sum = 0
  tracks.forEach(t => {
    if (completedTracks.value[t.id]) sum += t.rewardCash
  })
  return sum
})

const currentLevelTitle = computed(() => {
  const count = completedTracksCount.value
  if (count >= 4) return 'Master Quantitative Investor'
  if (count >= 2) return 'Savvy Portfolio Architect'
  if (count >= 1) return 'Aspiring Investor'
  return 'Novice Apprentice'
})

const badges = computed(() => [
  { id: 'b1', name: 'Compounding Guru', icon: '🌱', criteria: 'Complete Compound Growth lesson', unlocked: !!completedTracks.value['compounding'] },
  { id: 'b2', name: 'Shield of Diversification', icon: '🛡️', criteria: 'Master Asset Allocation', unlocked: !!completedTracks.value['diversification'] },
  { id: 'b3', name: 'Quant Analyst', icon: '⚡', criteria: 'Pass Sharpe & Beta mastery', unlocked: !!completedTracks.value['beta_sharpe'] },
  { id: 'b4', name: 'Risk Sentinel', icon: '👑', criteria: 'Complete all 4 academy tracks', unlocked: completedTracksCount.value === 4 }
])

function calculateCompounding(monthly) {
  const r = 0.10 / 12
  const n = 20 * 12
  const fv = monthly * ((Math.pow(1 + r, n) - 1) / r)
  return Math.round(fv)
}

function openLesson(track) {
  activeLesson.value = track
  selectedAnswer.value = null
  quizSubmitted.value = false
}

const isAnswerCorrect = computed(() => {
  if (!activeLesson.value || selectedAnswer.value === null) return false
  return selectedAnswer.value === activeLesson.value.quiz.correctIndex
})

function submitQuiz() {
  if (selectedAnswer.value === null) return
  quizSubmitted.value = true

  if (isAnswerCorrect.value) {
    const isFirstTime = !completedTracks.value[activeLesson.value.id]
    completedTracks.value[activeLesson.value.id] = true
    localStorage.setItem('investwise_completed_lessons', JSON.stringify(completedTracks.value))

    if (isFirstTime) {
      emit('reward-user', activeLesson.value.rewardCash)
    }
  }
}

function finishLesson() {
  activeLesson.value = null
}
</script>

<style scoped>
.academy-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.academy-hero {
  padding: 2rem;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.12));
  border-color: rgba(99, 102, 241, 0.3);
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.4);
  font-size: 0.75rem;
  font-weight: 800;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  margin-bottom: 0.75rem;
  letter-spacing: 0.05em;
}

.hero-title {
  font-size: 1.75rem;
  font-weight: 800;
  line-height: 1.25;
  margin-bottom: 0.5rem;
}

.hero-desc {
  font-size: 0.95rem;
  color: var(--text-muted);
  max-width: 700px;
  line-height: 1.5;
  margin-bottom: 1.5rem;
}

.progress-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 0.85rem;
}

.stat-pill {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--border-color);
  padding: 0.6rem 1rem;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
}

.pill-label {
  font-size: 0.7rem;
  color: var(--text-dim);
  text-transform: uppercase;
  font-weight: 600;
}

.pill-val {
  font-family: var(--font-heading);
  font-size: 1.1rem;
  font-weight: 700;
}

.text-accent { color: #a5b4fc; }
.text-success { color: var(--success); }

/* Active Lesson Card */
.active-lesson-card {
  padding: 2rem;
  border-color: var(--primary);
  background: rgba(15, 23, 42, 0.95);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.6);
}

@media (max-width: 640px) {
  .active-lesson-card {
    padding: 1.25rem;
  }
}

.lesson-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.25rem;
}

.lesson-title {
  font-size: 1.4rem;
  font-weight: 800;
  margin-top: 0.25rem;
}

.lesson-content-body {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.lesson-illustration {
  text-align: center;
  font-size: 3rem;
  background: rgba(255, 255, 255, 0.03);
  padding: 1rem;
  border-radius: 16px;
  border: 1px solid var(--border-color);
}

.concept-heading {
  font-size: 1.15rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: #c7d2fe;
}

.concept-p {
  font-size: 0.95rem;
  line-height: 1.55;
  color: #e2e8f0;
  margin-bottom: 1rem;
}

.mini-calc-box {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1.25rem;
}

.calc-header {
  font-size: 0.85rem;
  font-weight: 700;
  color: #a5b4fc;
  margin-bottom: 0.75rem;
}

.calc-slider-wrap {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  font-size: 0.85rem;
  margin-bottom: 0.75rem;
}

.calc-slider {
  accent-color: var(--primary);
}

.calc-result {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  font-size: 0.9rem;
}

.key-takeaways {
  background: rgba(245, 158, 11, 0.08);
  border-left: 3px solid var(--warning);
  padding: 0.85rem 1rem;
  border-radius: 0 10px 10px 0;
  font-size: 0.875rem;
}

.takeaway-title {
  font-weight: 700;
  color: #fbbf24;
  display: block;
  margin-bottom: 0.2rem;
}

.quiz-section {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 1.25rem;
}

.quiz-title {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: #a5b4fc;
}

.quiz-question {
  font-size: 0.95rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.quiz-options {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
  margin-bottom: 1.25rem;
}

.quiz-opt-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.85rem 1rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  color: var(--text-main);
  text-align: left;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.quiz-opt-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--border-active);
}

.opt-letter {
  font-weight: 700;
  background: rgba(255, 255, 255, 0.1);
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
}

.opt-selected {
  border-color: var(--primary);
  background: rgba(99, 102, 241, 0.15);
}

.opt-correct {
  border-color: var(--success) !important;
  background: rgba(16, 185, 129, 0.2) !important;
}

.opt-wrong {
  border-color: var(--danger) !important;
  background: rgba(244, 63, 94, 0.2) !important;
}

.quiz-feedback {
  padding: 1rem;
  border-radius: 10px;
  margin-bottom: 1.25rem;
}

.feedback-success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.35);
}

.feedback-fail {
  background: rgba(244, 63, 94, 0.15);
  border: 1px solid rgba(244, 63, 94, 0.35);
}

.feedback-title {
  font-weight: 700;
  font-size: 0.95rem;
  margin-bottom: 0.25rem;
}

.feedback-desc {
  font-size: 0.85rem;
  color: #e2e8f0;
}

.reward-notice {
  margin-top: 0.5rem;
  font-weight: 700;
  color: #34d399;
  font-size: 0.9rem;
}

/* Tracks Grid */
.tracks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.25rem;
}

.track-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s;
}

.track-card:hover {
  transform: translateY(-3px);
  border-color: var(--primary);
}

.track-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.track-icon-wrap {
  font-size: 2rem;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
}

.track-title {
  font-size: 1.15rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.track-summary {
  font-size: 0.85rem;
  color: var(--text-muted);
  line-height: 1.45;
  margin-bottom: 1.25rem;
  flex-grow: 1;
}

.track-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border-color);
}

.track-meta {
  font-size: 0.75rem;
  color: var(--text-dim);
}

/* Badges Showcase */
.badges-card {
  padding: 1.5rem;
}

.badges-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.badge-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border-color);
  padding: 1rem 0.75rem;
  border-radius: 12px;
  opacity: 0.45;
  filter: grayscale(80%);
  transition: all 0.3s;
}

.badge-item.badge-unlocked {
  opacity: 1;
  filter: none;
  background: rgba(99, 102, 241, 0.12);
  border-color: rgba(99, 102, 241, 0.4);
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.2);
}

.badge-icon-box {
  font-size: 2rem;
  margin-bottom: 0.35rem;
}

.badge-name {
  font-weight: 700;
  font-size: 0.85rem;
  margin-bottom: 0.2rem;
}

.badge-crit {
  font-size: 0.7rem;
  color: var(--text-dim);
}
</style>
