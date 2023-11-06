package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type key string

const (
	TxKey key = "tx"
)

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

// RollbackTx rollback transaction and log if error occurs.
func RollbackTx(ctx context.Context, tx pgx.Tx, log *zap.Logger, msg string, fields ...zap.Field) {
	err := tx.Rollback(ctx)
	if err == nil || errors.Is(err, pgx.ErrTxClosed) {
		return
	}

	log.With(zap.Error(err)).Error(msg, fields...)
}

type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

func BuildQuery(name string, sq Sqlizer) (Query, error) {
	raw, args, err := sq.ToSql()
	if err != nil {
		return Query{}, err
	}

	q := Query{
		Name:     name,
		QueryRaw: raw,
		Args:     args,
	}

	return q, nil
}
