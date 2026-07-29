package portfolio

import (
	"context"
	"database/sql"
	"time"
)

// Portfolio represents the user account cash balance and identification
type Portfolio struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	CashBalance float64   `json:"cash_balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Holding represents aggregates shares owned and average buy cost for a ticker
type Holding struct {
	ID              int       `json:"id"`
	PortfolioID     int       `json:"portfolio_id"`
	Ticker          string    `json:"ticker"`
	Shares          float64   `json:"shares"`
	AverageBuyPrice float64   `json:"average_buy_price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Transaction represents audit logs of all buys and sells
type Transaction struct {
	ID              int64     `json:"id"`
	PortfolioID     int       `json:"portfolio_id"`
	Ticker          string    `json:"ticker"`
	TransactionType string    `json:"transaction_type"` // "BUY" or "SELL"
	Shares          float64   `json:"shares"`
	Price           float64   `json:"price"`
	TransactionTime time.Time `json:"transaction_time"`
}

// PortfolioSummary represents real-time aggregated portfolio values & performance metrics
type PortfolioSummary struct {
	PortfolioID       int                  `json:"portfolio_id"`
	CashBalance       float64              `json:"cash_balance"`
	HoldingsValue     float64              `json:"holdings_value"`
	TotalValue        float64              `json:"total_value"`
	RealizedPnL       float64              `json:"realized_pnl"`
	UnrealizedPnL     float64              `json:"unrealized_pnl"`
	Holdings          []*HoldingDetail     `json:"holdings"`
	SectorAllocations map[string]float64   `json:"sector_allocations"` // Hardcoded mapping by sector
}

// HoldingDetail wraps standard holding with current valuation details
type HoldingDetail struct {
	Ticker          string  `json:"ticker"`
	Shares          float64 `json:"shares"`
	AverageBuyPrice float64 `json:"average_buy_price"`
	CurrentPrice    float64 `json:"current_price"`
	MarketValue     float64 `json:"market_value"`
	UnrealizedPnL   float64 `json:"unrealized_pnl"`
	AllocationPct   float64 `json:"allocation_pct"`
}

// PortfolioRepository describes interfaces for SQL portfolio storage
type PortfolioRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateUser(ctx context.Context, username, email, passwordHash string) (int, error)
	GetUserByEmail(ctx context.Context, email string) (id int, passwordHash string, err error)
	CreatePortfolio(ctx context.Context, userID int, initialCash float64) (int, error)
	GetPortfolioByUserID(ctx context.Context, userID int) (*Portfolio, error)
	GetPortfolioByIDForUpdate(ctx context.Context, tx *sql.Tx, portfolioID int) (*Portfolio, error)
	UpdatePortfolioCash(ctx context.Context, tx *sql.Tx, portfolioID int, cash float64) error
	GetHolding(ctx context.Context, portfolioID int, ticker string) (*Holding, error)
	GetHoldings(ctx context.Context, portfolioID int) ([]*Holding, error)
	UpsertHolding(ctx context.Context, tx *sql.Tx, holding *Holding) error
	DeleteHolding(ctx context.Context, tx *sql.Tx, portfolioID int, ticker string) error
	CreateTransaction(ctx context.Context, tx *sql.Tx, txModel *Transaction) error
	GetTransactions(ctx context.Context, portfolioID int) ([]*Transaction, error)
}
