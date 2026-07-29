package risk

import (
	"context"
	"fmt"

	"github.com/Tabhi109/investwise/internal/logger"
)

// PortfolioDataProvider abstracts fetching portfolio holdings and cash balance
type PortfolioDataProvider interface {
	GetHoldingsAndCash(ctx context.Context, portfolioID int) (holdings map[string]float64, cash float64, err error)
}

// MarketDataProvider abstracts fetching current prices and historical price series
type MarketDataProvider interface {
	GetPrice(ctx context.Context, ticker string) (float64, error)
	GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error)
}

// Service coordinates calculating portfolio analytics
type Service struct {
	portfolioProvider PortfolioDataProvider
	marketProvider    MarketDataProvider
	riskFreeRate      float64
}

// NewService instantiates a Risk Service
func NewService(p PortfolioDataProvider, m MarketDataProvider, riskFreeRate float64) *Service {
	return &Service{
		portfolioProvider: p,
		marketProvider:    m,
		riskFreeRate:      riskFreeRate,
	}
}

// CalculatePortfolioRisk computes full risk metrics report for a portfolio ID
func (s *Service) CalculatePortfolioRisk(ctx context.Context, portfolioID int) (*PortfolioRiskMetrics, error) {
	// 1. Fetch portfolio holdings and cash
	holdings, cash, err := s.portfolioProvider.GetHoldingsAndCash(ctx, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve portfolio data: %w", err)
	}

	// 2. Compute individual asset market values and total portfolio value
	assetValues := make(map[string]float64)
	totalAssetValue := 0.0

	for ticker, shares := range holdings {
		if shares <= 0 {
			continue
		}
		price, err := s.marketProvider.GetPrice(ctx, ticker)
		if err != nil {
			logger.Warn("Failed to fetch price for risk calculation, skipping holding", "ticker", ticker, "err", err)
			continue
		}
		val := shares * price
		assetValues[ticker] = val
		totalAssetValue += val
	}

	totalPortfolioValue := totalAssetValue + cash
	if totalPortfolioValue <= 0 {
		return &PortfolioRiskMetrics{
			AssetWeights:      make(map[string]float64),
			RiskContributions: make(map[string]float64),
		}, nil
	}

	// 3. Compute asset weights (weights sum to 1.0, including cash if cash > 0)
	weights := make(map[string]float64)
	for ticker, val := range assetValues {
		weights[ticker] = val / totalPortfolioValue
	}
	// Note: cash weight is 1.0 - sum(weights), it has 0 returns/volatility, handled implicitly by scaling

	if len(weights) == 0 {
		// All cash portfolio
		return &PortfolioRiskMetrics{
			ExpectedReturn:    0.0,
			Volatility:        0.0,
			Variance:          0.0,
			SharpeRatio:       0.0,
			Beta:              0.0,
			MaxDrawdown:       0.0,
			AssetWeights:      map[string]float64{"CASH": 1.0},
			RiskContributions: make(map[string]float64),
		}, nil
	}

	// 4. Retrieve historical returns for assets in portfolio + benchmark (SPY)
	historicalReturns := make(map[string][]float64)
	limit := 91 // We need 91 prices to get 90 daily returns

	for ticker := range weights {
		prices, err := s.marketProvider.GetHistoricalPrices(ctx, ticker, limit)
		if err != nil || len(prices) < 2 {
			return nil, fmt.Errorf("insufficient historical prices for %s: %w", ticker, err)
		}
		historicalReturns[ticker] = CalculateReturns(prices)
	}

	// Retrieve benchmark historical returns (SPY)
	benchPrices, err := s.marketProvider.GetHistoricalPrices(ctx, "SPY", limit)
	var benchmarkReturns []float64
	if err == nil && len(benchPrices) >= 2 {
		benchmarkReturns = CalculateReturns(benchPrices)
	} else {
		logger.Warn("Failed to retrieve benchmark returns for Beta calculation, benchmark returns will be empty")
	}

	// 5. Execute matrix computation via pure math layer
	metrics := ComputeRiskMetrics(totalPortfolioValue, weights, historicalReturns, benchmarkReturns, s.riskFreeRate)

	// Add CASH allocation explicitly to report if cash exists
	if cash > 0 {
		metrics.AssetWeights["CASH"] = cash / totalPortfolioValue
	}

	return metrics, nil
}
