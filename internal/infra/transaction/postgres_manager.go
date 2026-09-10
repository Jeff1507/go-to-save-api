package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

func Querier(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

type PostgresManager struct {
	pg *pgxpool.Pool
}

func NewPostgresManager(pg *pgxpool.Pool) *PostgresManager {
	return &PostgresManager{
		pg: pg,
	}
}

func (txm *PostgresManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := txm.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin manager transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		return fmt.Errorf("execute manager callback function: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit manager transaction: %w", err)
	}

	return nil
}
