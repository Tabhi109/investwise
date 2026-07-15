package market

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/Tabhi109/investwise/internal/logger"
	"github.com/Tabhi109/investwise/internal/websocket"
)

// Service coordinates fetching prices and running ticker updates
type Service struct {
	repo    MarketRepository
	tickers []string
	rand    *rand.Rand
}

// NewService instantiates a Market Service
func NewService(repo MarketRepository) *Service {
	return &Service{
		repo:    repo,
		tickers: []string{"AAPL", "MSFT", "GOOGL", "AMZN", "NVDA", "SPY"},
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetPrice fetches current price (checks Redis cache)
func (s *Service) GetPrice(ctx context.Context, ticker string) (*StockPrice, error) {
	return s.repo.GetPrice(ctx, ticker)
}

// GetTickers returns the list of simulated ticker symbols
func (s *Service) GetTickers() []string {
	return s.tickers
}

// GetHistoricalPrices retrieves chronological closing prices from database
func (s *Service) GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error) {
	return s.repo.GetHistoricalPrices(ctx, ticker, limit)
}

// SimulatePrices runs a loops generating stock prices, caching them, and broadcasting via WebSockets
func (s *Service) SimulatePrices(ctx context.Context, hub *websocket.Hub, interval time.Duration) {
	logger.Info("Starting market price simulation loop")

	basePrices := map[string]float64{
		"AAPL":  180.0,
		"MSFT":  400.0,
		"GOOGL": 150.0,
		"AMZN":  175.0,
		"NVDA":  800.0,
		"SPY":   500.0, // Benchmark S&P ETF
	}

	prices := make(map[string]float64)
	for _, ticker := range s.tickers {
		cached, err := s.repo.GetPrice(ctx, ticker)
		if err == nil && cached != nil {
			prices[ticker] = cached.Price
		} else {
			prices[ticker] = basePrices[ticker]
			_ = s.repo.SetPrice(ctx, &StockPrice{
				Ticker:    ticker,
				Price:     basePrices[ticker],
				Change:    0.0,
				Timestamp: time.Now(),
			})
		}
	}

	tickerChan := time.NewTicker(interval)
	defer tickerChan.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Market price simulation loop terminated")
			return
		case <-tickerChan.C:
			for _, ticker := range s.tickers {
				currPrice := prices[ticker]
				// Random walk from -0.5% to +0.52% (slight upward long drift)
				changePercent := (s.rand.Float64() - 0.49) * 0.01
				newPrice := currPrice * (1.0 + changePercent)
				prices[ticker] = newPrice

				sp := &StockPrice{
					Ticker:    ticker,
					Price:     newPrice,
					Change:    changePercent * 100.0,
					Timestamp: time.Now(),
				}

				// Cache updated price in Redis
				_ = s.repo.SetPrice(ctx, sp)

				// Format update and broadcast to subscribers of this ticker
				data, err := json.Marshal(map[string]any{
					"type":      "ticker",
					"ticker":    ticker,
					"price":     newPrice,
					"change":    changePercent * 100.0,
					"timestamp": sp.Timestamp.Format(time.RFC3339),
				})
				if err == nil {
					hub.Broadcast <- &websocket.BroadcastMessage{
						Ticker: ticker,
						Data:   data,
					}
				}
			}
		}
	}
}
