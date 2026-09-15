package repository

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/review"
	"github.com/Jeff1507/go-to-save-api/internal/infra/db"
	"github.com/Jeff1507/go-to-save-api/internal/infra/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pg *pgxpool.Pool
}

func NewPostgresRepo(pg *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		pg: pg,
	}
}

func (repo *PostgresRepo) Save(ctx context.Context, r *review.Review) error {
	tx, ok := transaction.Querier(ctx)
	if ok {
		return repo.save(ctx, tx, r)
	}

	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := repo.save(ctx, tx, r); err != nil {
		return fmt.Errorf("save with transaction: %w", err)
	}

	return tx.Commit(ctx)
}

func (repo *PostgresRepo) save(ctx context.Context, tx pgx.Tx, r *review.Review) error {
	rRow := toRow(r)

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO reviews(id, order_id, rating, comment, reviewed_at)
		VALUES ($1, $2, $3, $4, $5)
		`,
		rRow.ID, rRow.OrderID, rRow.Rating, rRow.Comment, rRow.ReviewedAt,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return review.ErrAlreadyReviewed
		}

		return fmt.Errorf("insert review: %w", err)
	}

	return nil
}
