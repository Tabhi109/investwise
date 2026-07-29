package portfolio

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tabhi109/investwise/internal/auth"
	"github.com/Tabhi109/investwise/internal/market"
	"github.com/Tabhi109/investwise/internal/risk"
	"github.com/gin-gonic/gin"
)

// Mock portfolio repository implementation
type mockPortfolioRepo struct {
	users        map[string]string // email -> pwdHash
	portfolios   map[int]*Portfolio // userID -> Portfolio
	transactions []*Transaction
	holdings     map[string]*Holding // ticker -> Holding
}

func newMockPortfolioRepo() *mockPortfolioRepo {
	return &mockPortfolioRepo{
		users:      make(map[string]string),
		portfolios: make(map[int]*Portfolio),
		holdings:   make(map[string]*Holding),
	}
}

func (m *mockPortfolioRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, nil // No-op for mock
}

func (m *mockPortfolioRepo) CreateUser(ctx context.Context, username, email, passwordHash string) (int, error) {
	if _, exists := m.users[email]; exists {
		return 0, errors.New("user already exists")
	}
	m.users[email] = passwordHash
	return 1, nil
}

func (m *mockPortfolioRepo) GetUserByEmail(ctx context.Context, email string) (int, string, error) {
	hash, exists := m.users[email]
	if !exists {
		return 0, "", sql.ErrNoRows
	}
	return 1, hash, nil
}

func (m *mockPortfolioRepo) CreatePortfolio(ctx context.Context, userID int, initialCash float64) (int, error) {
	m.portfolios[userID] = &Portfolio{
		ID:          100,
		UserID:      userID,
		CashBalance: initialCash,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return 100, nil
}

func (m *mockPortfolioRepo) GetPortfolioByUserID(ctx context.Context, userID int) (*Portfolio, error) {
	p, exists := m.portfolios[userID]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return p, nil
}

func (m *mockPortfolioRepo) GetPortfolioByIDForUpdate(ctx context.Context, tx *sql.Tx, portfolioID int) (*Portfolio, error) {
	return m.portfolios[1], nil // Default mocked portfolio
}

func (m *mockPortfolioRepo) UpdatePortfolioCash(ctx context.Context, tx *sql.Tx, portfolioID int, cash float64) error {
	for _, p := range m.portfolios {
		if p.ID == portfolioID {
			p.CashBalance = cash
			return nil
		}
	}
	return sql.ErrNoRows
}

func (m *mockPortfolioRepo) GetHolding(ctx context.Context, portfolioID int, ticker string) (*Holding, error) {
	h, exists := m.holdings[ticker]
	if !exists {
		return nil, nil
	}
	return h, nil
}

func (m *mockPortfolioRepo) GetHoldings(ctx context.Context, portfolioID int) ([]*Holding, error) {
	var list []*Holding
	for _, h := range m.holdings {
		list = append(list, h)
	}
	return list, nil
}

func (m *mockPortfolioRepo) UpsertHolding(ctx context.Context, tx *sql.Tx, holding *Holding) error {
	m.holdings[holding.Ticker] = holding
	return nil
}

func (m *mockPortfolioRepo) DeleteHolding(ctx context.Context, tx *sql.Tx, portfolioID int, ticker string) error {
	delete(m.holdings, ticker)
	return nil
}

func (m *mockPortfolioRepo) CreateTransaction(ctx context.Context, tx *sql.Tx, txModel *Transaction) error {
	txModel.ID = int64(len(m.transactions) + 1)
	txModel.TransactionTime = time.Now()
	m.transactions = append(m.transactions, txModel)
	return nil
}

func (m *mockPortfolioRepo) GetTransactions(ctx context.Context, portfolioID int) ([]*Transaction, error) {
	return m.transactions, nil
}

// Mock market repository for service mapping
type mockMarketRepo struct {
	prices map[string]*market.StockPrice
}

func (m *mockMarketRepo) GetPrice(ctx context.Context, ticker string) (*market.StockPrice, error) {
	p, exists := m.prices[ticker]
	if !exists {
		return nil, errors.New("price not found")
	}
	return p, nil
}

func (m *mockMarketRepo) SetPrice(ctx context.Context, price *market.StockPrice) error {
	m.prices[price.Ticker] = price
	return nil
}

func (m *mockMarketRepo) GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error) {
	return []float64{100.0, 101.0, 102.0, 103.0}, nil
}

func TestRegisterAndLoginIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Instantiate mock repos
	pRepo := newMockPortfolioRepo()
	mRepo := &mockMarketRepo{prices: make(map[string]*market.StockPrice)}

	// Setup Services
	mService := market.NewService(mRepo)
	pService := NewService(pRepo, mService)
	authService := auth.NewService("test-secret-key-1234567890", 1)
	
	// Set mock provider for risk engine
	riskService := risk.NewService(pService, &mockMarketProvider{mRepo}, 0.04)

	handler := NewHandler(pService, authService, riskService)

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// 1. Test registration
	regPayload := registerRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(regPayload)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	// 2. Test login
	loginPayload := loginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var loginResponse map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"].(string)
	if !ok || token == "" {
		t.Fatal("login response did not contain JWT token string")
	}
}

type mockMarketProvider struct {
	repo *mockMarketRepo
}

func (m *mockMarketProvider) GetPrice(ctx context.Context, ticker string) (float64, error) {
	p, err := m.repo.GetPrice(ctx, ticker)
	if err != nil {
		return 0, err
	}
	return p.Price, nil
}

func (m *mockMarketProvider) GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error) {
	return m.repo.GetHistoricalPrices(ctx, ticker, limit)
}
