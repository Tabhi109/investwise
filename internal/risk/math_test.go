package risk

import (
	"math"
	"testing"
)

func TestCalculateReturns(t *testing.T) {
	prices := []float64{100.0, 105.0, 102.9, 108.045}
	// returns: [0.05, -0.02, 0.05]
	expected := []float64{0.05, -0.02, 0.05}

	returns := CalculateReturns(prices)
	if len(returns) != len(expected) {
		t.Fatalf("expected returns length %d, got %d", len(expected), len(returns))
	}

	for i := range expected {
		if math.Abs(returns[i]-expected[i]) > 1e-6 {
			t.Errorf("at index %d: expected return %f, got %f", i, expected[i], returns[i])
		}
	}
}

func TestCalculateMeanAndVariance(t *testing.T) {
	values := []float64{0.05, -0.02, 0.05}
	expectedMean := 0.02666666666666667
	expectedVar := 0.0016333333333333334 // sample variance formula

	mean := CalculateMean(values)
	if math.Abs(mean-expectedMean) > 1e-6 {
		t.Errorf("expected mean %f, got %f", expectedMean, mean)
	}

	variance := CalculateVariance(values, mean)
	if math.Abs(variance-expectedVar) > 1e-6 {
		t.Errorf("expected variance %f, got %f", expectedVar, variance)
	}
}

func TestComputeRiskMetrics(t *testing.T) {
	// Construct simulated price histories for 2 assets: AAPL and MSFT
	// returns of AAPL: [0.01, -0.01, 0.02, -0.02, 0.03]
	// returns of MSFT: [0.02,  0.01, 0.01, -0.01, 0.02]
	historicalReturns := map[string][]float64{
		"AAPL": {0.01, -0.01, 0.02, -0.02, 0.03},
		"MSFT": {0.02, 0.01, 0.01, -0.01, 0.02},
	}

	// Benchmark (SPY) returns: [0.015, 0.0, 0.01, -0.015, 0.025]
	benchmarkReturns := []float64{0.015, 0.0, 0.01, -0.015, 0.025}

	weights := map[string]float64{
		"AAPL": 0.6,
		"MSFT": 0.4,
	}

	portfolioValue := 100000.0
	riskFreeRate := 0.04

	metrics := ComputeRiskMetrics(portfolioValue, weights, historicalReturns, benchmarkReturns, riskFreeRate)

	// Verify math identity: The sum of risk contributions must equal the annualized portfolio volatility
	var riskContributionSum float64
	for _, contribution := range metrics.RiskContributions {
		riskContributionSum += contribution
	}

	if math.Abs(riskContributionSum-metrics.Volatility) > 1e-7 {
		t.Errorf("Mathematical Identity Violation: sum of asset risk contributions (%f) must equal annualized portfolio volatility (%f)", riskContributionSum, metrics.Volatility)
	}

	// Volatility should be positive
	if metrics.Volatility <= 0 {
		t.Errorf("expected positive annualized volatility, got %f", metrics.Volatility)
	}

	// HHI for 60/40 should be 0.6^2 + 0.4^2 = 0.36 + 0.16 = 0.52
	expectedHHI := 0.52
	if math.Abs(metrics.HHI-expectedHHI) > 1e-6 {
		t.Errorf("expected HHI %f, got %f", expectedHHI, metrics.HHI)
	}

	// Parametric VaR should be positive
	if metrics.ParametricVaR95 <= 0 || metrics.ParametricVaR99 <= 0 {
		t.Errorf("expected positive parametric VaR, got VaR95=%f, VaR99=%f", metrics.ParametricVaR95, metrics.ParametricVaR99)
	}

	// Historical VaR should be positive
	if metrics.HistoricalVaR95 < 0 || metrics.HistoricalVaR99 < 0 {
		t.Errorf("expected non-negative historical VaR, got VaR95=%f, VaR99=%f", metrics.HistoricalVaR95, metrics.HistoricalVaR99)
	}
}
