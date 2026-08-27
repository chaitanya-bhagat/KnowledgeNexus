package logger

import (
	"fmt"
	"strings"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type LoggerConfig struct {
// 	Environment string
// 	Level       string
// }

func NewLogger(cfg config.LoggerConfig) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var zapCfg zap.Config
	if strings.EqualFold(cfg.Environment, "development") {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapCfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger %w", err)
	}
	return logger, nil
}

func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "", "info":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level %q", level)
	}
}
