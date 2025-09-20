package main

import (
	"log/slog"
	"os"
	"strconv"
)

type config struct {
	addr   string
	logger *slog.Logger
	db     struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func loadConfig() config {
	var cfg config

	// Dale, let's setup the server configuration
	// Default port is 4000, but you can change it if you want, no drama
	cfg.addr = getEnv("SERVER_ADDR", ":4000")

	// Aca va la config de la DB, re importante esto eh!
	// If you mess this up, everything goes to la mierda
	cfg.db.dsn = getEnv("DB_DSN", "appuser:appusersecret@tcp(localhost:3306)/classifiersdb")
	
	// More conservative connection pool settings to avoid overwhelming the database
	cfg.db.maxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	cfg.db.maxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 25)
	cfg.db.maxIdleTime = getEnv("DB_MAX_IDLE_TIME", "15m")

	return cfg
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
