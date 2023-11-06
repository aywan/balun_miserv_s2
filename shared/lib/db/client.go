package db

import (
	"context"
	"fmt"

	"github.com/aywan/balun_miserv_s2/shared/lib/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	pgx_zap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

type PgClient struct {
	cfg  Config
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewDB(cfg Config, log *zap.Logger) *PgClient {
	return &PgClient{cfg: cfg, log: log}
}

func (p *PgClient) Start(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"user=%s password='%s' host=%s port=%s dbname=%s sslmode=disable",
		p.cfg.User,
		p.cfg.Pass,
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.Name,
	)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	if p.cfg.Trace {
		var logLevel tracelog.LogLevel
		logLevel, err = tracelog.LogLevelFromString(p.cfg.LogLevel)
		if err != nil {
			return err
		}
		poolCfg.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgx_zap.NewLogger(p.log),
			LogLevel: logLevel,
		}
	}

	p.pool, err = pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return err
	}

	err = p.pool.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PgClient) ScanOneContext(ctx context.Context, dest interface{}, q Query) error {
	row, err := p.QueryContext(ctx, q)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *PgClient) ScanAllContext(ctx context.Context, dest interface{}, q Query) error {
	rows, err := p.QueryContext(ctx, q)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *PgClient) ExecContext(ctx context.Context, q Query) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, q.Args...)
	}

	return p.pool.Exec(ctx, q.QueryRaw, q.Args...)
}

func (p *PgClient) QueryContext(ctx context.Context, q Query) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, q.Args...)
	}

	return p.pool.Query(ctx, q.QueryRaw, q.Args...)
}

func (p *PgClient) QueryRowContext(ctx context.Context, q Query) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, q.Args...)
	}

	return p.pool.QueryRow(ctx, q.QueryRaw, q.Args...)
}

func (p *PgClient) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}

func (p *PgClient) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *PgClient) Close() {
	p.pool.Close()
}

func (p *PgClient) ReadCommitted(ctx context.Context, f TxHandler) (err error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}
	log := logger.CtxLogDef(ctx, p.log)
	defer func() {
		if r := recover(); r != nil {
			log.With(zap.Any("recover", r)).Error("recover from panic")
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	tx, err = p.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	defer RollbackTx(ctx, tx, log, "rollback transaction")

	ctx = MakeContextTx(ctx, tx)

	err = f(ctx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
