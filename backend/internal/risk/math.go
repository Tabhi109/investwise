package risk

import (
	"math"
	"sort"
)

// CalculateReturns converts a series of prices to daily percentage returns.
// Returns list is of length N-1.
func CalculateReturns(prices []float64) []float64 {
	if len(prices) < 2 {
		return nil
	}
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if prices[i-1] == 0 {
			returns[i-1] = 0
			continue
		}
		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}
	return returns
}

// CalculateMean returns the arithmetic mean of a slice.
func CalculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// CalculateVariance computes the sample variance of a slice.
func CalculateVariance(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0
	}
	sumSq := 0.0
	for _, v := range values {
		diff := v - mean
		sumSq += diff * diff
	}
	return sumSq / float64(len(values)-1)
}

// CalculateCovariance computes the sample covariance between two slices of equal length.
func CalculateCovariance(x, y []float64, meanX, meanY float64) float64 {
	n := len(x)
	if n < 2 || n != len(y) {
		return 0
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += (x[i] - meanX) * (y[i] - meanY)
	}
	return sum / float64(n-1)
}

// ComputeRiskMetrics executes the quantitative matrix calculations for the portfolio.
func ComputeRiskMetrics(
	portfolioValue float64,
	weights map[string]float64,
	historicalReturns map[string][]float64,
	benchmarkReturns []float64,
	riskFreeRate float64,
) *PortfolioRiskMetrics {
	// 1. Setup list of tickers
	var tickers []string
	for ticker := range weights {
		tickers = append(tickers, ticker)
	}

	n := len(tickers)
	if n == 0 {
		return &PortfolioRiskMetrics{}
	}

	// 2. Compute mean daily return and variance (volatility) for each asset
	means := make(map[string]float64)
	vols := make(map[string]float64)
	minReturnLen := -1

	for _, ticker := range tickers {
		ret := historicalReturns[ticker]
		if len(ret) == 0 {
			continue
		}
		if minReturnLen == -1 || len(ret) < minReturnLen {
			minReturnLen = len(ret)
		}
		mean := CalculateMean(ret)
		means[ticker] = mean
		vols[ticker] = math.Sqrt(CalculateVariance(ret, mean))
	}

	if minReturnLen <= 0 {
		return &PortfolioRiskMetrics{}
	}

	// Truncate returns to align dates exactly (using the shortest available length)
	alignedReturns := make(map[string][]float64)
	for _, ticker := range tickers {
		ret := historicalReturns[ticker]
		alignedReturns[ticker] = ret[len(ret)-minReturnLen:]
	}

	var alignedBenchmark []float64
	if len(benchmarkReturns) >= minReturnLen {
		alignedBenchmark = benchmarkReturns[len(benchmarkReturns)-minReturnLen:]
	} else if len(benchmarkReturns) > 0 {
		minReturnLen = len(benchmarkReturns)
		alignedBenchmark = benchmarkReturns
		for _, ticker := range tickers {
			ret := alignedReturns[ticker]
			alignedReturns[ticker] = ret[len(ret)-minReturnLen:]
		}
	}

	// 3. Compute Covariance Matrix (daily)
	covMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		covMatrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		t1 := tickers[i]
		for j := i; j < n; j++ {
			t2 := tickers[j]
			cov := CalculateCovariance(alignedReturns[t1], alignedReturns[t2], means[t1], means[t2])
			covMatrix[i][j] = cov
			covMatrix[j][i] = cov
		}
	}

	// 4. Expected Portfolio Return (annualized)
	portfolioExpectedReturnDaily := 0.0
	hhi := 0.0
	weightedVolSum := 0.0
	for _, ticker := range tickers {
		w := weights[ticker]
		portfolioExpectedReturnDaily += w * means[ticker]
		hhi += w * w
		weightedVolSum += w * vols[ticker]
	}
	portfolioExpectedReturnAnnual := portfolioExpectedReturnDaily * 252.0

	// 5. Portfolio Variance & Volatility
	portfolioVarianceDaily := 0.0
	for i := 0; i < n; i++ {
		t1 := tickers[i]
		for j := 0; j < n; j++ {
			t2 := tickers[j]
			portfolioVarianceDaily += weights[t1] * weights[t2] * covMatrix[i][j]
		}
	}
	portfolioVolDaily := math.Sqrt(portfolioVarianceDaily)
	portfolioVolAnnual := portfolioVolDaily * math.Sqrt(252.0)

	// 6. Sharpe Ratio
	sharpe := 0.0
	if portfolioVolAnnual > 0 {
		sharpe = (portfolioExpectedReturnAnnual - riskFreeRate) / portfolioVolAnnual
	}

	// 7. Portfolio daily returns series (historical scenario paths)
	portfolioReturnsHistory := make([]float64, minReturnLen)
	for t := 0; t < minReturnLen; t++ {
		dailyRet := 0.0
		for _, ticker := range tickers {
			dailyRet += weights[ticker] * alignedReturns[ticker][t]
		}
		portfolioReturnsHistory[t] = dailyRet
	}

	// 8. Beta calculation
	beta := 1.0
	if len(alignedBenchmark) == minReturnLen && minReturnLen > 1 {
		meanBench := CalculateMean(alignedBenchmark)
		varBench := CalculateVariance(alignedBenchmark, meanBench)
		if varBench > 0 {
			meanPort := CalculateMean(portfolioReturnsHistory)
			covPortBench := CalculateCovariance(portfolioReturnsHistory, alignedBenchmark, meanPort, meanBench)
			beta = covPortBench / varBench
		}
	}

	// 9. Maximum Drawdown (MDD) calculation
	maxDrawdown := 0.0
	portfolioValSim := 100.0
	peakVal := 100.0
	for _, dailyRet := range portfolioReturnsHistory {
		portfolioValSim = portfolioValSim * (1.0 + dailyRet)
		if portfolioValSim > peakVal {
			peakVal = portfolioValSim
		}
		dd := (portfolioValSim - peakVal) / peakVal
		if dd < maxDrawdown {
			maxDrawdown = dd
		}
	}

	// 10. Value at Risk (VaR)
	// Parametric VaR: Z * Vol - Expected Return
	z95 := 1.64485
	z99 := 2.32635
	paramVaR95DailyPercent := (z95 * portfolioVolDaily) - portfolioExpectedReturnDaily
	paramVaR99DailyPercent := (z99 * portfolioVolDaily) - portfolioExpectedReturnDaily

	if paramVaR95DailyPercent < 0 {
		paramVaR95DailyPercent = 0
	}
	if paramVaR99DailyPercent < 0 {
		paramVaR99DailyPercent = 0
	}

	// Historical VaR: Percentile threshold of sorted historical returns
	sortedReturns := make([]float64, len(portfolioReturnsHistory))
	copy(sortedReturns, portfolioReturnsHistory)
	sort.Float64s(sortedReturns)

	histVaR95DailyPercent := 0.0
	histVaR99DailyPercent := 0.0
	if len(sortedReturns) > 0 {
		idx95 := int(math.Floor(float64(len(sortedReturns)) * 0.05))
		idx99 := int(math.Floor(float64(len(sortedReturns)) * 0.01))
		if idx95 < 0 {
			idx95 = 0
		}
		if idx99 < 0 {
			idx99 = 0
		}
		// VaR is expressed as positive loss
		histVaR95DailyPercent = -sortedReturns[idx95]
		histVaR99DailyPercent = -sortedReturns[idx99]

		if histVaR95DailyPercent < 0 {
			histVaR95DailyPercent = 0
		}
		if histVaR99DailyPercent < 0 {
			histVaR99DailyPercent = 0
		}
	}

	// Convert percentages to absolute cash terms
	parametricVaR95 := portfolioValue * paramVaR95DailyPercent
	parametricVaR99 := portfolioValue * paramVaR99DailyPercent
	historicalVaR95 := portfolioValue * histVaR95DailyPercent
	historicalVaR99 := portfolioValue * histVaR99DailyPercent

	// 11. Diversification Ratio
	diversificationRatio := 1.0
	if portfolioVolDaily > 0 {
		diversificationRatio = weightedVolSum / portfolioVolDaily
	}

	// 12. Risk Contribution per Asset
	riskContributions := make(map[string]float64)
	if portfolioVolDaily > 0 {
		for i, ticker := range tickers {
			// (Covariance matrix * weights vector) index i
			covSum := 0.0
			for j := 0; j < n; j++ {
				covSum += covMatrix[i][j] * weights[tickers[j]]
			}
			// Marginal Contribution = w_i * Cov(R_i, R_p) / sigma_p
			marginalContribution := weights[ticker] * (covSum / portfolioVolDaily)
			// Annualize it to present in consistent annual volatility units
			riskContributions[ticker] = marginalContribution * math.Sqrt(252.0)
		}
	}

	return &PortfolioRiskMetrics{
		ExpectedReturn:       portfolioExpectedReturnAnnual,
		Volatility:           portfolioVolAnnual,
		Variance:             portfolioVarianceDaily,
		SharpeRatio:          sharpe,
		Beta:                 beta,
		MaxDrawdown:          maxDrawdown,
		ParametricVaR95:      parametricVaR95,
		ParametricVaR99:      parametricVaR99,
		HistoricalVaR95:      historicalVaR95,
		HistoricalVaR99:      historicalVaR99,
		DiversificationRatio: diversificationRatio,
		HHI:                  hhi,
		RiskContributions:    riskContributions,
		AssetWeights:         weights,
	}
}
