package db

import (
	"context"
	"fmt"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/config"
	pgx_zap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

// New creates a new pool connection to the database.`
func New(ctx context.Context, cfg config.DbConfig, log *zap.Logger) (*pgxpool.Pool, error) {
	var err error

	dsn := fmt.Sprintf(
		"user=%s password='%s' host=%s port=%s dbname=%s sslmode=disable",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	if cfg.Trace {
		var logLevel tracelog.LogLevel
		logLevel, err = tracelog.LogLevelFromString(cfg.LogLevel)
		if err != nil {
			return nil, err
		}
		poolCfg.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgx_zap.NewLogger(log),
			LogLevel: logLevel,
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
