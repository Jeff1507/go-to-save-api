package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
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

func (repo *PostgresRepo) Save(ctx context.Context, b *bag.Bag) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	r := toRow(b)

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO bags(
			id, 
			store_id, 
			title, 
			description, 
			stock, 
			original_price, 
			final_price, 
			status,
			pickup_start, 
			pickup_end, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
		r.ID,
		r.StoreID,
		r.Title,
		r.Description,
		r.Stock,
		r.OriginalPrice,
		r.FinalPrice,
		r.Status,
		r.PickupWindow.Start,
		r.PickupWindow.End,
		r.CreatedAt,
		r.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert bag: %w", err)
	}

	batch := &pgx.Batch{}
	for _, item := range r.Items {
		batch.Queue(
			`
			INSERT INTO bag_items(id, bag_id, name, image_key, quantity)
			VALUES ($1, $2, $3, $4, $5)
			`,
			item.ID, r.ID, item.Name, item.ImageKey, item.Quantity,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range r.Items {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("insert bag item: %w", err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, b *bag.Bag) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	r := toRow(b)

	_, err = tx.Exec(
		ctx,
		`
		UPDATE bags
		SET 
			title = $1, 
			description = $2, 
			original_price = $3, 
			final_price = $4, 
			status = $5, 
			pickup_start = $6,
			pickup_end = $7, 
			updated_at = $8
		WHERE id = $9
		`,
		r.Title,
		r.Description,
		r.OriginalPrice,
		r.FinalPrice,
		r.Status,
		r.PickupWindow.Start,
		r.PickupWindow.End,
		r.UpdatedAt,
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("update bag: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM bag_items WHERE bag_id = $1
		`,
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("delete bag items: %w", err)
	}

	batch := &pgx.Batch{}
	for _, item := range r.Items {
		batch.Queue(
			`
			INSERT INTO bag_items(id, bag_id, name, image_key, quantity)
			VALUES ($1, $2, $3, $4, $5)
			`,
			item.ID, r.ID, item.Name, item.ImageKey, item.Quantity,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range r.Items {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("insert bag item: %w", err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*bag.Bag, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT
			id, 
			store_id, 
			title, 
			description, 
			stock, 
			original_price, 
			final_price, 
			status,
			pickup_start, 
			pickup_end, 
			created_at, 
			updated_at
		FROM bags
		WHERE id = $1
		`,
		id,
	)

	var r bagRow
	if err := row.Scan(
		&r.ID,
		&r.StoreID,
		&r.Title,
		&r.Description,
		&r.Stock,
		&r.OriginalPrice,
		&r.FinalPrice,
		&r.Status,
		&r.PickupWindow.Start,
		&r.PickupWindow.End,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, bag.ErrNotFound
		}

		return nil, fmt.Errorf("scan bag: %w", err)
	}

	rows, err := repo.pg.Query(
		ctx,
		`
		SELECT id, name, image_key, quantity
		FROM bag_items
		WHERE bag_id = $1
		`,
		r.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("select bag items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i itemRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ImageKey,
			&i.Quantity,
		); err != nil {
			return nil, fmt.Errorf("scan bag item: %w", err)
		}

		r.Items = append(r.Items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return toEntity(r), nil
}
