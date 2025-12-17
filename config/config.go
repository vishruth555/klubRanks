package config

import "time"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port    string
	Log     string
	Counter int
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

// Single exported config instance
var AppConfig = Config{
	Server: ServerConfig{
		Port:    "8080",
		Log:     "debug",
		Counter: 1,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
	},
	Database: DatabaseConfig{
		Driver: "sqlite3",
		DSN:    "klubranks.db",
	},
	JWT: JWTConfig{
		Secret: "super-secret-key",
		Expiry: 24 * time.Hour,
	},
}
