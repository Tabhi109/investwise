package risk

// PortfolioRiskMetrics holds computed portfolio risk analytics
type PortfolioRiskMetrics struct {
	ExpectedReturn       float64            `json:"expected_return"`       // Annualized expected return
	Volatility           float64            `json:"volatility"`            // Annualized volatility
	Variance             float64            `json:"variance"`              // Daily portfolio variance
	SharpeRatio          float64            `json:"sharpe_ratio"`          // Sharpe ratio relative to risk free rate
	Beta                 float64            `json:"beta"`                  // Volatility beta relative to benchmark index (SPY)
	MaxDrawdown          float64            `json:"max_drawdown"`          // Maximum drawdown in historical simulations
	ParametricVaR95      float64            `json:"parametric_var_95"`      // Parametric 1-day 95% VaR (in cash terms)
	ParametricVaR99      float64            `json:"parametric_var_99"`      // Parametric 1-day 99% VaR (in cash terms)
	HistoricalVaR95      float64            `json:"historical_var_95"`      // Historical 1-day 95% VaR (in cash terms)
	HistoricalVaR99      float64            `json:"historical_var_99"`      // Historical 1-day 99% VaR (in cash terms)
	DiversificationRatio float64            `json:"diversification_ratio"` // Diversification ratio (weighted avg vol / port vol)
	HHI                  float64            `json:"hhi"`                   // Herfindahl-Hirschman concentration score
	RiskContributions    map[string]float64 `json:"risk_contributions"`    // Marginal risk contribution per asset (in volatility units)
	AssetWeights         map[string]float64 `json:"asset_weights"`         // Asset allocation weights
}
