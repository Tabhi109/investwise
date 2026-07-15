package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tabhi109/investwise/internal/auth"
	"github.com/Tabhi109/investwise/internal/config"
	"github.com/Tabhi109/investwise/internal/database"
	"github.com/Tabhi109/investwise/internal/logger"
	"github.com/Tabhi109/investwise/internal/market"
	"github.com/Tabhi109/investwise/internal/portfolio"
	"github.com/Tabhi109/investwise/internal/redis"
	"github.com/Tabhi109/investwise/internal/risk"
	"github.com/Tabhi109/investwise/internal/worker"
	"github.com/Tabhi109/investwise/internal/websocket"
	"github.com/gin-gonic/gin"
)

// marketRiskAdapter adapts market.Service to satisfy risk.MarketDataProvider interface
type marketRiskAdapter struct {
	marketService *market.Service
}

func (a *marketRiskAdapter) GetPrice(ctx context.Context, ticker string) (float64, error) {
	p, err := a.marketService.GetPrice(ctx, ticker)
	if err != nil {
		return 0, err
	}
	return p.Price, nil
}

func (a *marketRiskAdapter) GetHistoricalPrices(ctx context.Context, ticker string, limit int) ([]float64, error) {
	return a.marketService.GetHistoricalPrices(ctx, ticker, limit)
}

// Application represents the monolith composition root.
type Application struct {
	Config      *config.Config
	router      *gin.Engine
	postgres    *database.Postgres
	redis       *redis.Redis
	workerPool  *worker.WorkerPool
	wsHub       *websocket.Hub
	authService *auth.Service
}

// NewApplication coordinates DI setup
func NewApplication() *Application {
	// 1. Load config and initialize logger
	cfg := config.Load()
	logger.Init(cfg.Env, os.Stdout)
	logger.Info("Bootstrapping InvestWise Monolith")

	// Set Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return &Application{
		Config: cfg,
		router: router,
	}
}

// Run performs setup connection connections, starts worker pools, starts router and handles graceful shutdown
func (a *Application) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Connect Postgres
	postgres, err := database.Connect(ctx, a.Config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.postgres = postgres
	logger.Info("Connected to PostgreSQL database")

	// 2. Connect Redis
	redisClient, err := redis.Connect(ctx, a.Config.RedisURL)
	if err != nil {
		_ = postgres.Close()
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	a.redis = redisClient
	logger.Info("Connected to Redis cache")

	// 3. Setup WebSocket Hub and Worker Pool
	a.wsHub = websocket.NewHub()
	a.workerPool = worker.NewWorkerPool(a.Config.WorkerCount)
	a.authService = auth.NewService(a.Config.JWTSecret, a.Config.JWTTTLHours)

	// 4. Construct Business Components (DI Layer)
	marketRepo := market.NewRepository(a.postgres, a.redis)
	marketService := market.NewService(marketRepo)
	marketHandler := market.NewHandler(marketService)

	portfolioRepo := portfolio.NewRepository(a.postgres)
	portfolioService := portfolio.NewService(portfolioRepo, marketService)

	// Wire risk engine adapters
	mAdapter := &marketRiskAdapter{marketService: marketService}
	riskService := risk.NewService(portfolioService, mAdapter, a.Config.RiskFreeRate)
	portfolioHandler := portfolio.NewHandler(portfolioService, a.authService, riskService)

	// Register Routes
	a.setupRoutes(portfolioHandler, marketHandler, a.wsHub, a.authService)

	// 5. Start Background Workers
	go a.wsHub.Run()
	a.workerPool.Start()

	// Submit market simulation loop to worker pool
	interval := time.Duration(a.Config.MarketUpdateIntervalMs) * time.Millisecond
	a.workerPool.Submit(func(wCtx context.Context) {
		marketService.SimulatePrices(wCtx, a.wsHub, interval)
	})

	// 6. Bind HTTP Server
	server := &http.Server{
		Addr:    ":" + a.Config.Port,
		Handler: a.router,
	}

	serverError := make(chan error, 1)
	go func() {
		logger.Info("HTTP server listening", "port", a.Config.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	// 7. Await Termination Signal (SIGINT, SIGTERM)
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverError:
		logger.Error("HTTP server failed, starting shutdown", err)
	case sig := <-shutdownSignal:
		logger.Info("Termination signal received, initiating graceful shutdown", "signal", sig.String())
	}

	// 8. Graceful Shutdown Sequence
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Stop HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown failed", err)
	} else {
		logger.Info("HTTP server stopped gracefully")
	}

	// Stop WebSocket connections
	a.wsHub.Close()
	logger.Info("WebSocket Hub terminated")

	// Stop Background Worker Pool
	a.workerPool.Stop()

	// Close database connections
	if err := a.postgres.Close(); err != nil {
		logger.Error("Failed to close PostgreSQL pool connection", err)
	}
	if err := a.redis.Close(); err != nil {
		logger.Error("Failed to close Redis connection", err)
	}

	logger.Info("Monolith application shutdown complete")
	return nil
}

func (a *Application) setupRoutes(
	portHandler *portfolio.Handler,
	mktHandler *market.Handler,
	wsHub *websocket.Hub,
	authService *auth.Service,
) {
	// Health check endpoint
	a.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
	})

	// Public auth and market routes
	a.router.POST("/register", portHandler.Register)
	a.router.POST("/login", portHandler.Login)
	a.router.GET("/tickers", mktHandler.GetTickers)
	a.router.GET("/market/price/:ticker", mktHandler.GetPrice)

	// Realtime feed upgrades (handled by WebSockets)
	a.router.GET("/ws", websocket.HandleWS(wsHub))

	// Protected routes group (requires JWT header)
	protected := a.router.Group("/")
	protected.Use(authService.AuthMiddleware())
	{
		protected.GET("/portfolio", portHandler.GetSummary)
		protected.POST("/trade", portHandler.Trade)
		protected.GET("/portfolio/transactions", portHandler.GetTransactions)
		protected.GET("/portfolio/risk", portHandler.GetRiskMetrics)
	}
}
