package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"subscriptions/internal/api"
	"subscriptions/internal/logger"
	"subscriptions/internal/storage/postgresClient"
)

// Config defines configuration parameters for the notification-service application,
// including HTTP server setting, SMTP/PostreSQL/Redis credentials, logger optional and calculate timeouts.
type Config struct {
	HttpServer api.HttpServer
	Postgres   postgresClient.Config
	Logger     logger.Config
}

// New loads the configuration from the specified file path and initializes computed timeout values.
// Returns a fully filled Config instance or an error if loading fails.
func New(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
