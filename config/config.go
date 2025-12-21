package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port            string
	Log             string
	Counter         int
	CoolDownMinutes int
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Log:             getEnv("LOG_LEVEL", "info"),
			Counter:         1,
			CoolDownMinutes: 1,
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite3"),
			DSN:    getEnv("DB_DSN", "klubranks.db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "dev-secret"),
			Expiry: 48 * time.Hour,
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
