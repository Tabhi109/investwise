package market

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tabhi109/investwise/internal/database"
	"github.com/Tabhi109/investwise/internal/redis"
)

// Repository implements MarketRepository using Postgres and Redis cache
type Repository struct {
	db    *database.Postgres
	redis *redis.Redis
}

// NewRepository initializes the Market repository adapter
func NewRepository(db *database.Postgres, redisClient *redis.Redis) MarketRepository {
	return &Repository{
		db:    db,
		redis: redisClient,
	}
}

// GetPrice retrieves the cached price from Redis
func (r *Repository) GetPrice(ctx context.Context, ticker string) (*StockPrice, error) {
	key := fmt.Sprintf("market:price:%s", ticker)
	data, err := r.redis.Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("cached price not found for %s: %w", ticker, err)
	}

	var price StockPrice
	if err := json.Unmarshal(data, &price); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached price: %w", err)
	}

	return &price, nil
}

// SetPrice caches the current price in Redis with a TTL of 1 hour (refreshed constantly by feed)
func (r *Repository) SetPrice(ctx context.Context, price *StockPrice) error {
	key := fmt.Sprintf("market:price:%s", price.Ticker)
	data, err := json.Marshal(price)
	if err != nil {
		return fmt.Errorf("failed to marshal price: %w", err)
	}

	// 1-hour expiration represents a fallback TTL
	return r.redis.Client.Set(ctx, key, data, time.Hour).Err()
}

// GetHistoricalPrices queries database for last N prices ordered chronologically
func (r *Repository) GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error) {
	// Nested query retrieves the latest N dates, then sorts them ascending by price_date (chronological)
	query := `
		SELECT close_price 
		FROM (
			SELECT close_price, price_date 
			FROM asset_prices 
			WHERE ticker = $1 
			ORDER BY price_date DESC 
			LIMIT $2
		) sub 
		ORDER BY price_date ASC`

	rows, err := r.db.DB.QueryContext(ctx, query, ticker, limit)
	if err != nil {
		return nil, fmt.Errorf("historical prices query failed: %w", err)
	}
	defer rows.Close()

	var prices []float64
	for rows.Next() {
		var price float64
		if err := rows.Scan(&price); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		prices = append(prices, price)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return prices, nil
}
