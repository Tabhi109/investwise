package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tabhi109/investwise/internal/auth"
	"github.com/Tabhi109/investwise/internal/market"
)

// Service coordinates business logic for portfolios, trades, and user accounts
type Service struct {
	repo          PortfolioRepository
	marketService *market.Service
}

// NewService instantiates a Portfolio Service
func NewService(repo PortfolioRepository, m *market.Service) *Service {
	return &Service{
		repo:          repo,
		marketService: m,
	}
}

// RegisterUser hashes password, creates database user record, and provisions default portfolio with $100k cash
func (s *Service) RegisterUser(ctx context.Context, username, email, password string) (int, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.repo.CreateUser(ctx, username, email, hash)
	if err != nil {
		return 0, err
	}

	// Default starting simulated cash: $100,000.00
	_, err = s.repo.CreatePortfolio(ctx, userID, 100000.0000)
	if err != nil {
		return 0, fmt.Errorf("failed to provision user portfolio: %w", err)
	}

	return userID, nil
}

// LoginUser validates email/password and returns user ID if valid
func (s *Service) LoginUser(ctx context.Context, email, password string) (int, error) {
	userID, hash, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	if !auth.CheckPasswordHash(password, hash) {
		return 0, errors.New("invalid credentials")
	}

	return userID, nil
}

// GetHoldingsAndCash implements risk.PortfolioDataProvider interface
func (s *Service) GetHoldingsAndCash(ctx context.Context, portfolioID int) (map[string]float64, float64, error) {
	// Standard query bypasses GORM/GORM transactions.
	// Since we are read-only for risk, we don't lock.
	p, err := s.repo.GetPortfolioByUserID(ctx, portfolioID) // portfolio ID matches user ID 1-to-1 in current setup
	if err != nil {
		// Try fetching by portfolio ID directly if needed (fallback/support user ID matches)
		// For simplicity, users have exactly one portfolio. So we query user ID as portfolio ID.
		p = &Portfolio{ID: portfolioID}
		// Let's resolve portfolio by user ID
		p, err = s.repo.GetPortfolioByUserID(ctx, portfolioID)
		if err != nil {
			return nil, 0, err
		}
	}

	holdings, err := s.repo.GetHoldings(ctx, p.ID)
	if err != nil {
		return nil, 0, err
	}

	holdingsMap := make(map[string]float64)
	for _, h := range holdings {
		if h.Shares > 0 {
			holdingsMap[h.Ticker] = h.Shares
		}
	}

	return holdingsMap, p.CashBalance, nil
}

// GetPortfolioSummary compiles valuation aggregates and calculates dynamic realized/unrealized performance
func (s *Service) GetPortfolioSummary(ctx context.Context, userID int) (*PortfolioSummary, error) {
	p, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load portfolio: %w", err)
	}

	holdings, err := s.repo.GetHoldings(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load holdings: %w", err)
	}

	// Map sectors
	sectors := map[string]string{
		"AAPL":  "Technology",
		"MSFT":  "Technology",
		"GOOGL": "Technology",
		"AMZN":  "Consumer Discretionary",
		"NVDA":  "Technology",
		"SPY":   "Index",
	}

	var holdingsValue float64
	var unrealizedPnL float64
	var details []*HoldingDetail
	sectorAlloc := make(map[string]float64)

	// Fetch current market prices and calculate holdings valuation
	for _, h := range holdings {
		quote, err := s.marketService.GetPrice(ctx, h.Ticker)
		currentPrice := h.AverageBuyPrice // Fallback if no quote yet
		if err == nil && quote != nil {
			currentPrice = quote.Price
		}

		mktVal := h.Shares * currentPrice
		unrealPnL := h.Shares * (currentPrice - h.AverageBuyPrice)

		holdingsValue += mktVal
		unrealizedPnL += unrealPnL

		details = append(details, &HoldingDetail{
			Ticker:          h.Ticker,
			Shares:          h.Shares,
			AverageBuyPrice: h.AverageBuyPrice,
			CurrentPrice:    currentPrice,
			MarketValue:     mktVal,
			UnrealizedPnL:   unrealPnL,
		})

		sec := sectors[h.Ticker]
		if sec == "" {
			sec = "Other"
		}
		sectorAlloc[sec] += mktVal
	}

	totalValue := holdingsValue + p.CashBalance

	// Add allocation percentages
	for _, d := range details {
		if totalValue > 0 {
			d.AllocationPct = d.MarketValue / totalValue
		}
	}

	// Compute sector allocation percentages
	for k, v := range sectorAlloc {
		if totalValue > 0 {
			sectorAlloc[k] = v / totalValue
		}
	}
	if p.CashBalance > 0 && totalValue > 0 {
		sectorAlloc["Cash"] = p.CashBalance / totalValue
	}

	// Compute Realized PnL dynamically from transactions audit history
	txs, err := s.repo.GetTransactions(ctx, p.ID)
	realizedPnL := 0.0
	if err == nil {
		realizedPnL = s.calculateRealizedPnL(txs)
	}

	return &PortfolioSummary{
		PortfolioID:       p.ID,
		CashBalance:       p.CashBalance,
		HoldingsValue:     holdingsValue,
		TotalValue:        totalValue,
		RealizedPnL:       realizedPnL,
		UnrealizedPnL:     unrealizedPnL,
		Holdings:          details,
		SectorAllocations: sectorAlloc,
	}, nil
}

// calculateRealizedPnL computes realized profit/losses dynamically using FIFO/Average Buy Cost rules
func (s *Service) calculateRealizedPnL(txs []*Transaction) float64 {
	// Process chronological history of transactions (which is returned descending, so we reverse it)
	if len(txs) == 0 {
		return 0
	}

	// Reverse transaction slice to process oldest first
	n := len(txs)
	chronoTxs := make([]*Transaction, n)
	for i, t := range txs {
		chronoTxs[n-1-i] = t
	}

	realizedPnL := 0.0
	// Keep track of average cost simulation to match DB state updates
	sharesMap := make(map[string]float64)
	avgCostMap := make(map[string]float64)

	for _, t := range chronoTxs {
		if t.TransactionType == "BUY" {
			currentShares := sharesMap[t.Ticker]
			currentAvg := avgCostMap[t.Ticker]

			newShares := currentShares + t.Shares
			newAvg := ((currentShares * currentAvg) + (t.Shares * t.Price)) / newShares

			sharesMap[t.Ticker] = newShares
			avgCostMap[t.Ticker] = newAvg
		} else if t.TransactionType == "SELL" {
			currentShares := sharesMap[t.Ticker]
			currentAvg := avgCostMap[t.Ticker]

			if currentShares >= t.Shares {
				// Realized Profit = shares sold * (sell price - average buy cost)
				realizedPnL += t.Shares * (t.Price - currentAvg)
				sharesMap[t.Ticker] = currentShares - t.Shares
			}
		}
	}

	return realizedPnL
}

// ExecuteTrade processes BUY/SELL transactions in database using transaction isolation and row-level locks
func (s *Service) ExecuteTrade(ctx context.Context, userID int, ticker string, tradeType string, shares float64) error {
	if tradeType != "BUY" && tradeType != "SELL" {
		return errors.New("invalid trade type: must be BUY or SELL")
	}
	if shares <= 0 {
		return errors.New("shares must be positive")
	}

	// Retrieve simulated current price
	quote, err := s.marketService.GetPrice(ctx, ticker)
	if err != nil || quote == nil {
		return fmt.Errorf("failed to fetch current market price for %s: %w", ticker, err)
	}
	price := quote.Price

	// Start Database Transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin database transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Lock portfolio row
	p, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user portfolio: %w", err)
	}

	p, err = s.repo.GetPortfolioByIDForUpdate(ctx, tx, p.ID)
	if err != nil {
		return fmt.Errorf("failed to lock portfolio: %w", err)
	}

	holding, err := s.repo.GetHolding(ctx, p.ID, ticker)
	if err != nil {
		return fmt.Errorf("failed to check holdings: %w", err)
	}

	if tradeType == "BUY" {
		cost := shares * price
		if p.CashBalance < cost {
			return fmt.Errorf("insufficient cash balance: have $%f, need $%f", p.CashBalance, cost)
		}

		newCash := p.CashBalance - cost
		if err := s.repo.UpdatePortfolioCash(ctx, tx, p.ID, newCash); err != nil {
			return fmt.Errorf("failed to update cash: %w", err)
		}

		var newShares float64
		var newAvgPrice float64

		if holding != nil {
			newShares = holding.Shares + shares
			newAvgPrice = ((holding.Shares * holding.AverageBuyPrice) + (shares * price)) / newShares
		} else {
			newShares = shares
			newAvgPrice = price
		}

		err = s.repo.UpsertHolding(ctx, tx, &Holding{
			PortfolioID:     p.ID,
			Ticker:          ticker,
			Shares:          newShares,
			AverageBuyPrice: newAvgPrice,
		})
		if err != nil {
			return fmt.Errorf("failed to update holdings: %w", err)
		}

	} else if tradeType == "SELL" {
		if holding == nil || holding.Shares < shares {
			hasShares := 0.0
			if holding != nil {
				hasShares = holding.Shares
			}
			return fmt.Errorf("insufficient shares to sell: have %f shares, trying to sell %f", hasShares, shares)
		}

		revenue := shares * price
		newCash := p.CashBalance + revenue
		if err := s.repo.UpdatePortfolioCash(ctx, tx, p.ID, newCash); err != nil {
			return fmt.Errorf("failed to update cash: %w", err)
		}

		if holding.Shares == shares {
			// Closed position completely
			err = s.repo.DeleteHolding(ctx, tx, p.ID, ticker)
		} else {
			newShares := holding.Shares - shares
			err = s.repo.UpsertHolding(ctx, tx, &Holding{
				PortfolioID:     p.ID,
				Ticker:          ticker,
				Shares:          newShares,
				AverageBuyPrice: holding.AverageBuyPrice, // Buy average remains constant
			})
		}
		if err != nil {
			return fmt.Errorf("failed to update holdings: %w", err)
		}
	}

	// Write transaction row
	err = s.repo.CreateTransaction(ctx, tx, &Transaction{
		PortfolioID:     p.ID,
		Ticker:          ticker,
		TransactionType: tradeType,
		Shares:          shares,
		Price:           price,
	})
	if err != nil {
		return fmt.Errorf("failed to log audit transaction: %w", err)
	}

	// Commit Transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

// GetTransactions retrieves the trade logs list
func (s *Service) GetTransactions(ctx context.Context, userID int) ([]*Transaction, error) {
	p, err := s.repo.GetPortfolioByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTransactions(ctx, p.ID)
}
