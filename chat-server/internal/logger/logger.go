package logger

import (
	"github.com/aywan/balun_miserv_s2/chat-server/internal/config"
	"go.uber.org/zap"
)

// MustNewLogger initializes new logger.
func MustNewLogger(mode config.Mode) *zap.Logger {
	var logger *zap.Logger
	var err error

	switch mode {
	case config.DevMode:
		logger, err = zap.NewDevelopment()

	default:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return logger
}
