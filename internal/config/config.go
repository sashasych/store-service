package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// HTTP holds HTTP server configuration.
type HTTP struct {
	Addr string `envconfig:"HTTP_ADDR" default:":8080"`
}

// Postgres holds connection settings for PostgreSQL.
type Postgres struct {
	DSN      string `envconfig:"POSTGRES_DSN" default:"postgres://postgres:postgres@localhost:5432/store?sslmode=disable"`
	MaxConns int32  `envconfig:"POSTGRES_MAX_CONNS" default:"10"`
}

// Config is the root configuration structure populated from environment variables.
type Config struct {
	HTTP            HTTP
	Postgres        Postgres
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" default:"10s"`
	LogLevel        string        `envconfig:"LOG_LEVEL" default:"info"`
}

// Load reads configuration values from the environment.
func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
