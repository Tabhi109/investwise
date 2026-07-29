package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/0001_init.up.sql
var initSQL string

// Postgres wrapping the database/sql connection pool
type Postgres struct {
	DB *sql.DB
}

// Connect initializes the PostgreSQL connection pool using pgx driver.
func Connect(ctx context.Context, connStr string) (*Postgres, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	// Configure reasonable connection pool limits for a financial app
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connectivity
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	p := &Postgres{DB: db}

	// Auto-run migrations and seed data
	if err := p.InitSchemaAndSeed(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to run migrations/seed: %w", err)
	}

	return p, nil
}

// Close gracefully tears down the database connection pool
func (p *Postgres) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// InitSchemaAndSeed runs the embedded schema DDL and seeds initial mock prices
func (p *Postgres) InitSchemaAndSeed(ctx context.Context) error {
	// Execute schema migration
	_, err := p.DB.ExecContext(ctx, initSQL)
	if err != nil {
		return fmt.Errorf("migration exec failed: %w", err)
	}

	// Check if asset_prices is empty
	var count int
	err = p.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM asset_prices").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check asset_prices count: %w", err)
	}

	if count > 0 {
		// Already seeded
		return nil
	}

	// Seed historical data
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	tickers := map[string]float64{
		"AAPL":  180.0,
		"MSFT":  400.0,
		"GOOGL": 150.0,
		"AMZN":  175.0,
		"NVDA":  800.0,
		"SPY":   500.0, // Benchmark
	}

	// Seed 90 days of history
	now := time.Now()
	r := rand.New(rand.NewSource(42)) // Deterministic seed

	for ticker, basePrice := range tickers {
		currentPrice := basePrice
		dayCount := 0
		for i := 120; i >= 0; i-- { // 120 calendar days back
			date := now.AddDate(0, 0, -i)
			// Skip weekends
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				continue
			}

			// Random walk with a slight upward drift
			changePercent := (r.Float64() - 0.48) * 0.02 // drift +0.04% average
			currentPrice = currentPrice * (1.0 + changePercent)

			_, err = tx.ExecContext(ctx,
				"INSERT INTO asset_prices (ticker, price_date, close_price) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
				ticker, date.Format("2006-01-02"), currentPrice,
			)
			if err != nil {
				return fmt.Errorf("failed to insert historical price for %s: %w", ticker, err)
			}
			dayCount++
			if dayCount >= 90 { // Ensure exactly 90 trading days
				break
			}
		}
	}

	return tx.Commit()
}
