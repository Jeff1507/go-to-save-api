package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
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

func (repo *PostgresRepo) Save(ctx context.Context, p *product.Product) error {
	r := toRow(p)

	_, err := repo.pg.Exec(
		ctx,
		`
		INSERT INTO products(id, store_id, name, image_key, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		r.ID, r.StoreID, r.Name, r.ImageKey, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return product.ErrAlreadyExists
		}

		return fmt.Errorf("insert product: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, p *product.Product) error {
	r := toRow(p)

	_, err := repo.pg.Exec(
		ctx,
		`
		UPDATE products
		SET name = $1, image_key = $2, updated_at = $3
		WHERE id = $4
		`,
		r.Name, r.ImageKey, r.UpdatedAt, r.ID,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return product.ErrAlreadyExists
		}

		return fmt.Errorf("update product: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.pg.Exec(
		ctx,
		`
		DELETE FROM products WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete product: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*product.Product, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT id, store_id, name, image_key, created_at, updated_at
		FROM products
		WHERE id = $1
		`,
		id,
	)

	var r productRow
	if err := row.Scan(
		&r.ID,
		&r.StoreID,
		&r.Name,
		&r.ImageKey,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product.ErrNotFound
		}

		return nil, fmt.Errorf("scan product: %v", err)
	}

	return toEntity(r), nil
}

func (repo *PostgresRepo) List(ctx context.Context, ids []string) ([]*product.Product, error) {
	rows, err := repo.pg.Query(
		ctx,
		`
		SELECT id, store_id, name, image_key, created_at, updated_at
		FROM products
		WHERE id = ANY($1)
		`,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf("select products: %v", err)
	}
	defer rows.Close()

	products := make([]*product.Product, 0, len(ids))
	for rows.Next() {
		var r productRow
		if err := rows.Scan(
			&r.ID,
			&r.StoreID,
			&r.Name,
			&r.ImageKey,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %v", err)
		}

		products = append(products, toEntity(r))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %v", err)
	}

	return products, nil
}
