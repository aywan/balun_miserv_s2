package app

import (
	"context"

	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"github.com/aywan/balun_miserv_s2/shared/lib/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDb(lc fx.Lifecycle, log *zap.Logger, config *db.Config) *db.PgClient {
	dbInstance := db.NewDB(
		*config,
		logger.WithPart(log, "db"),
	)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dbInstance.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			dbInstance.Close()
			return nil
		},
	})
	return dbInstance
}
