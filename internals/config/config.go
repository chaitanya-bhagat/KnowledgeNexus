package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Logger     LoggerConfig
	Postgresql PostgresConfig
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	server, err := loadServerConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load server config: %w", err)
	}

	logger, err := loadLoggerConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load server config: %w", err)
	}

	postgresdb, err := loadPostgresConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load postgresql config: %w", err)
	}

	return Config{
		Server:     server,
		Logger:     logger,
		Postgresql: postgresdb,
	}, nil
}
