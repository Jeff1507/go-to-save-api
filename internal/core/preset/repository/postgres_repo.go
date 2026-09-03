package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
	"github.com/Jeff1507/go-to-save-api/internal/infra/db"
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

func (repo *PostgresRepo) Save(ctx context.Context, p *preset.Preset) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		INSERT INTO pickup_presets(id, store_id, name, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5)
		`,
		p.ID(), p.StoreID(), p.Name(), p.Start().String(), p.End().String(),
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return preset.ErrAlreadyExists
		}

		return fmt.Errorf("insert preset: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, p *preset.Preset) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		UPDATE pickup_presets
		SET name = $1, start_time = $2, end_time = $3
		WHERE id = $4
		`,
		p.Name(), p.Start().String(), p.End().String(), p.ID(),
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return preset.ErrAlreadyExists
		}

		return fmt.Errorf("update preset: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		DELETE FROM pickup_presets WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete preset: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*preset.Preset, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT id, store_id, name, start_time, end_time
		FROM pickup_presets
		WHERE id = $1
		`,
		id,
	)

	var r presetRow
	if err := row.Scan(
		&r.ID,
		&r.StoreID,
		&r.Name,
		&r.Start,
		&r.End,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, preset.ErrNotFound
		}

		return nil, fmt.Errorf("scan preset: %w", err)
	}

	return toEntity(r), nil
}
