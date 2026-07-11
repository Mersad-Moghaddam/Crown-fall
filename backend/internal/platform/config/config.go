package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddress     string
	LogLevel        string
	DatabaseURL     string
	RedisURL        string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPAddress:     value("CROWNFALL_HTTP_ADDRESS", ":8080"),
		LogLevel:        value("CROWNFALL_LOG_LEVEL", "info"),
		DatabaseURL:     os.Getenv("CROWNFALL_DATABASE_URL"),
		RedisURL:        os.Getenv("CROWNFALL_REDIS_URL"),
		ShutdownTimeout: duration("CROWNFALL_SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func value(key, fallback string) string {
	if result := os.Getenv(key); result != "" {
		return result
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	result, err := time.ParseDuration(os.Getenv(key))
	if err != nil || result <= 0 {
		return fallback
	}
	return result
}
