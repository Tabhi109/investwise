package market

import (
	"context"
	"time"
)

// StockPrice represents the current simulated price of a security ticker
type StockPrice struct {
	Ticker    string    `json:"ticker"`
	Price     float64   `json:"price"`
	Change    float64   `json:"change"` // percentage change from baseline or last tick
	Timestamp time.Time `json:"timestamp"`
}

// MarketRepository defines the persistence interface for stock prices (realtime in cache, historical in DB)
type MarketRepository interface {
	GetPrice(ctx context.Context, ticker string) (*StockPrice, error)
	SetPrice(ctx context.Context, price *StockPrice) error
	GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error)
}
