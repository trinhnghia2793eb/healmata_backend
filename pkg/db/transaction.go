package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct {
}

type SQLTxManager struct {
	pool *pgxpool.Pool
}

func NewSQLTxManager(pool *pgxpool.Pool) *SQLTxManager {
	return &SQLTxManager{
		pool: pool,
	}
}

func (tm *SQLTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Begin transaction
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	// rollback if error / not commit
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// apply tx into ctx
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	// execute callback logic (call to Repository)
	if err := fn(ctxWithTx); err != nil {
		return err
	}

	// callback successfully --> commit
	return tx.Commit(ctx)
}

// get Tx from context (for repository)
func GetTxFromContext(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}