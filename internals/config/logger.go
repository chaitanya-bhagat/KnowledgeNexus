package config

const (
	LOG_ENVIRONMENT = "LOG_ENVIRONMENT"
	LOG_LEVEL       = "LOG_LEVEL"
)

type LoggerConfig struct {
	Environment string
	Level       string
}

func loadLoggerConfig() (LoggerConfig, error) {
	return LoggerConfig{
		Environment: envOrDefault(LOG_ENVIRONMENT, "development"),
		Level:       envOrDefault(LOG_LEVEL, "info"),
	}, nil
}
