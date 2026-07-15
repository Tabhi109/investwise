package portfolio

import (
	"net/http"

	"github.com/Tabhi109/investwise/internal/auth"
	"github.com/Tabhi109/investwise/internal/risk"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	authService *auth.Service
	riskService *risk.Service
}

func NewHandler(s *Service, a *auth.Service, r *risk.Service) *Handler {
	return &Handler{
		service:     s,
		authService: a,
		riskService: r,
	}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type tradeRequest struct {
	Ticker string  `json:"ticker" binding:"required"`
	Type   string  `json:"type" binding:"required,oneof=BUY SELL"`
	Shares float64 `json:"shares" binding:"required,gt=0"`
}

// Register registers a new user and provisions portfolio
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.RegisterUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userID})
}

// Login authenticates credentials and issues token
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := h.authService.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": userID})
}

// GetSummary retrieves full portfolio performance details
func (h *Handler) GetSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	summary, err := h.service.GetPortfolioSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Trade processes Buy/Sell actions
func (h *Handler) Trade(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var req tradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ExecuteTrade(c.Request.Context(), userID, req.Ticker, req.Type, req.Shares)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trade executed successfully"})
}

// GetTransactions fetches historic execution orders
func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	txs, err := h.service.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txs)
}

// GetRiskMetrics triggers calculation from the Quantitative Risk Engine
func (h *Handler) GetRiskMetrics(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	// User ID maps to portfolio ID 1-to-1 in current database setup.
	metrics, err := h.riskService.CalculatePortfolioRisk(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
