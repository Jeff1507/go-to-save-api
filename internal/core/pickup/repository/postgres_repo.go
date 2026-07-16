package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/infra/db"
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

func (repo *PostgresRepo) Save(ctx context.Context, t *pickup.Template) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		INSERT INTO pickup_templates(id, store_id, name, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5)
		`,
		t.ID(), t.StoreID(), t.Name(), t.Start().String(), t.End().String(),
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return pickup.ErrTemplateAlreadyExists
		}

		return fmt.Errorf("insert pickup template: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, t *pickup.Template) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		UPDATE pickup_templates
		SET name = $1, start_time = $2, end_time = $3
		WHERE id = $4
		`,
		t.Name(), t.Start().String(), t.End().String(), t.ID(),
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return pickup.ErrTemplateAlreadyExists
		}

		return fmt.Errorf("update pickup template: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		DELETE FROM pickup_templates WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete pickup template: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*pickup.Template, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT id, store_id, name, start_time, end_time
		FROM pickup_templates
		WHERE id = $1
		`,
		id,
	)

	var r pickupTemplateRow
	if err := row.Scan(
		&r.ID,
		&r.StoreID,
		&r.Name,
		&r.Start,
		&r.End,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pickup.ErrTemplateNotFound
		}

		return nil, fmt.Errorf("scan pickup template: %v", err)
	}

	return toEntity(r), nil
}
