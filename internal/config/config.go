package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration parameters.
type Config struct {
	Env                    string
	Port                   string
	DatabaseURL            string
	RedisURL               string
	JWTSecret              string
	JWTTTLHours            int
	WorkerCount            int
	RiskFreeRate           float64
	MarketUpdateIntervalMs int
}

// Load reads configuration from environment variables, fallback to defaults.
func Load() *Config {
	return &Config{
		Env:                    getEnv("APP_ENV", "development"),
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/investwise?sslmode=disable"),
		RedisURL:               getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:              getEnv("JWT_SECRET", "investwise-dev-jwt-secret-key-change-in-prod"),
		JWTTTLHours:            getEnvAsInt("JWT_TTL_HOURS", 24),
		WorkerCount:            getEnvAsInt("WORKER_COUNT", 5),
		RiskFreeRate:           getEnvAsFloat("RISK_FREE_RATE", 0.04),
		MarketUpdateIntervalMs: getEnvAsInt("MARKET_UPDATE_INTERVAL_MS", 2000),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
}
