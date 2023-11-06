package logger

import (
	"context"

	"go.uber.org/zap"
)

type Mode string

const (
	// DevMode - means application in development mode.
	DevMode = Mode("dev")
	// ProdMode - means application in production mode.
	ProdMode = Mode("prod")
)

// MustNew initializes new logger.
func MustNew(mode Mode) *zap.Logger {
	logger, err := New(mode)

	if err != nil {
		panic(err)
	}

	return logger
}

// New initializes new logger.
func New(mode Mode) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch mode {
	case DevMode:
		logger, err = zap.NewDevelopment()

	default:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}

func WithPart(log *zap.Logger, part string) *zap.Logger {
	return log.With(zap.String("part", part))
}

type ctxKey string

var logCtxKey = ctxKey("ctx_log")

func CtxWithLog(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey, log)
}

func CtxLog(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(logCtxKey).(*zap.Logger); ok {
		return log
	}

	return nil
}

func MustCtxLog(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(logCtxKey).(*zap.Logger); ok {
		return log
	}

	panic("no logger in context")
}

func CtxLogDef(ctx context.Context, defLog *zap.Logger) *zap.Logger {
	if log, ok := ctx.Value(logCtxKey).(*zap.Logger); ok {
		return log
	}
	return defLog
}
