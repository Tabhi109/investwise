package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tabhi109/investwise/internal/database"
)

// Repository implements PortfolioRepository interface using database/sql
type Repository struct {
	postgres *database.Postgres
}

// NewRepository initializes Portfolio repository adapter
func NewRepository(postgres *database.Postgres) PortfolioRepository {
	return &Repository{
		postgres: postgres,
	}
}

// BeginTx starts a database transaction
func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.postgres.DB.BeginTx(ctx, nil)
}

// CreateUser registers a new user
func (r *Repository) CreateUser(ctx context.Context, username, email, passwordHash string) (int, error) {
	query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := r.postgres.DB.QueryRowContext(ctx, query, username, email, passwordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user in database: %w", err)
	}
	return id, nil
}

// GetUserByEmail retrieves credentials for validation
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (int, string, error) {
	query := "SELECT id, password_hash FROM users WHERE email = $1"
	var id int
	var hash string
	err := r.postgres.DB.QueryRowContext(ctx, query, email).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", fmt.Errorf("user not found: %w", err)
		}
		return 0, "", fmt.Errorf("failed to fetch user: %w", err)
	}
	return id, hash, nil
}

// CreatePortfolio provisions a new cash balance account for user
func (r *Repository) CreatePortfolio(ctx context.Context, userID int, initialCash float64) (int, error) {
	query := "INSERT INTO portfolios (user_id, cash_balance) VALUES ($1, $2) RETURNING id"
	var id int
	err := r.postgres.DB.QueryRowContext(ctx, query, userID, initialCash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create portfolio: %w", err)
	}
	return id, nil
}

// GetPortfolioByUserID returns user portfolio
func (r *Repository) GetPortfolioByUserID(ctx context.Context, userID int) (*Portfolio, error) {
	query := "SELECT id, user_id, cash_balance, created_at, updated_at FROM portfolios WHERE user_id = $1"
	p := &Portfolio{}
	err := r.postgres.DB.QueryRowContext(ctx, query, userID).Scan(&p.ID, &p.UserID, &p.CashBalance, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetPortfolioByIDForUpdate locks the row for transaction updates
func (r *Repository) GetPortfolioByIDForUpdate(ctx context.Context, tx *sql.Tx, portfolioID int) (*Portfolio, error) {
	query := "SELECT id, user_id, cash_balance, created_at, updated_at FROM portfolios WHERE id = $1 FOR UPDATE"
	p := &Portfolio{}
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, portfolioID).Scan(&p.ID, &p.UserID, &p.CashBalance, &p.CreatedAt, &p.UpdatedAt)
	} else {
		err = r.postgres.DB.QueryRowContext(ctx, query, portfolioID).Scan(&p.ID, &p.UserID, &p.CashBalance, &p.CreatedAt, &p.UpdatedAt)
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// UpdatePortfolioCash applies new cash balance
func (r *Repository) UpdatePortfolioCash(ctx context.Context, tx *sql.Tx, portfolioID int, cash float64) error {
	query := "UPDATE portfolios SET cash_balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, cash, portfolioID)
	} else {
		_, err = r.postgres.DB.ExecContext(ctx, query, cash, portfolioID)
	}
	return err
}

// GetHolding retrieves specific stock holding details
func (r *Repository) GetHolding(ctx context.Context, portfolioID int, ticker string) (*Holding, error) {
	query := "SELECT id, portfolio_id, ticker, shares, average_buy_price, created_at, updated_at FROM holdings WHERE portfolio_id = $1 AND ticker = $2"
	h := &Holding{}
	err := r.postgres.DB.QueryRowContext(ctx, query, portfolioID, ticker).Scan(&h.ID, &h.PortfolioID, &h.Ticker, &h.Shares, &h.AverageBuyPrice, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No holding exists
		}
		return nil, err
	}
	return h, nil
}

// GetHoldings gets list of all non-empty holdings
func (r *Repository) GetHoldings(ctx context.Context, portfolioID int) ([]*Holding, error) {
	query := "SELECT id, portfolio_id, ticker, shares, average_buy_price, created_at, updated_at FROM holdings WHERE portfolio_id = $1 AND shares > 0"
	rows, err := r.postgres.DB.QueryContext(ctx, query, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holdings []*Holding
	for rows.Next() {
		h := &Holding{}
		err := rows.Scan(&h.ID, &h.PortfolioID, &h.Ticker, &h.Shares, &h.AverageBuyPrice, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		holdings = append(holdings, h)
	}
	return holdings, nil
}

// UpsertHolding creates or modifies aggregated holding positions
func (r *Repository) UpsertHolding(ctx context.Context, tx *sql.Tx, holding *Holding) error {
	query := `
		INSERT INTO holdings (portfolio_id, ticker, shares, average_buy_price) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (portfolio_id, ticker) 
		DO UPDATE SET 
			shares = EXCLUDED.shares, 
			average_buy_price = EXCLUDED.average_buy_price, 
			updated_at = CURRENT_TIMESTAMP`
	
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, holding.PortfolioID, holding.Ticker, holding.Shares, holding.AverageBuyPrice)
	} else {
		_, err = r.postgres.DB.ExecContext(ctx, query, holding.PortfolioID, holding.Ticker, holding.Shares, holding.AverageBuyPrice)
	}
	return err
}

// DeleteHolding removes holding record (typically when shares reach 0)
func (r *Repository) DeleteHolding(ctx context.Context, tx *sql.Tx, portfolioID int, ticker string) error {
	query := "DELETE FROM holdings WHERE portfolio_id = $1 AND ticker = $2"
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, portfolioID, ticker)
	} else {
		_, err = r.postgres.DB.ExecContext(ctx, query, portfolioID, ticker)
	}
	return err
}

// CreateTransaction appends transaction row to audit trail
func (r *Repository) CreateTransaction(ctx context.Context, tx *sql.Tx, txModel *Transaction) error {
	query := "INSERT INTO transactions (portfolio_id, ticker, transaction_type, shares, price) VALUES ($1, $2, $3, $4, $5) RETURNING id, transaction_time"
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, txModel.PortfolioID, txModel.Ticker, txModel.TransactionType, txModel.Shares, txModel.Price).Scan(&txModel.ID, &txModel.TransactionTime)
	} else {
		err = r.postgres.DB.QueryRowContext(ctx, query, txModel.PortfolioID, txModel.Ticker, txModel.TransactionType, txModel.Shares, txModel.Price).Scan(&txModel.ID, &txModel.TransactionTime)
	}
	return err
}

// GetTransactions retrieves the user transaction history log
func (r *Repository) GetTransactions(ctx context.Context, portfolioID int) ([]*Transaction, error) {
	query := "SELECT id, portfolio_id, ticker, transaction_type, shares, price, transaction_time FROM transactions WHERE portfolio_id = $1 ORDER BY transaction_time DESC"
	rows, err := r.postgres.DB.QueryContext(ctx, query, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*Transaction
	for rows.Next() {
		txModel := &Transaction{}
		err := rows.Scan(&txModel.ID, &txModel.PortfolioID, &txModel.Ticker, &txModel.TransactionType, &txModel.Shares, &txModel.Price, &txModel.TransactionTime)
		if err != nil {
			return nil, err
		}
		txs = append(txs, txModel)
	}
	return txs, nil
}
