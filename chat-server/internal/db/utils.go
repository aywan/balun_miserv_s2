package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// RollbackTx rollback transaction and log if error occurs.
func RollbackTx(ctx context.Context, tx pgx.Tx, log *zap.Logger, msg string, fields ...zap.Field) {
	err := tx.Rollback(ctx)
	if err == nil || errors.Is(err, pgx.ErrTxClosed) {
		return
	}

	log.With(zap.Error(err)).Error(msg, fields...)
}
