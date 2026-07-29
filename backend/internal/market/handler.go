package market

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// GetTickers returns lists of active simulation tickers
func (h *Handler) GetTickers(c *gin.Context) {
	tickers := h.service.GetTickers()
	c.JSON(http.StatusOK, gin.H{"tickers": tickers})
}

// GetPrice retrieves cached quote for specific asset ticker
func (h *Handler) GetPrice(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticker is required"})
		return
	}

	price, err := h.service.GetPrice(c.Request.Context(), ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Price data not found"})
		return
	}

	c.JSON(http.StatusOK, price)
}
