package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"info"`
}

type HTTPConfig struct {
	Port    int    `env:"PORT" envDefault:"8080"`
	GinMode string `env:"GIN_MODE" envDefault:"release"`
}

type DatabaseConfig struct {
	URL             string        `env:"DATABASE_URL,required"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}
