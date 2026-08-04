package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
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

func (repo *PostgresRepo) Save(ctx context.Context, c *customer.Customer) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	r := toRow(c)

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO users(id, email, password_hash, verified_at)
		VALUES ($1, $2, $3, $4)
		`,
		r.ID, r.Email, r.PasswordHash, r.VerifiedAt,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return customer.ErrEmailAlreadyInUse
		}

		return fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO customers(id, first_name, last_name, picture_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		r.ID, r.FirstName, r.LastName, r.PictureKey, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, c *customer.Customer) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	r := toRow(c)

	_, err = tx.Exec(
		ctx,
		`
		UPDATE users SET email = $1, password_hash = $2, verified_at = $3
		WHERE id = $4
		`,
		r.Email, r.PasswordHash, r.VerifiedAt, r.ID,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return customer.ErrEmailAlreadyInUse
		}

		return fmt.Errorf("update user: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE customers SET first_name = $1, last_name = $2, picture_key = $3, updated_at = $4
		WHERE id = $5
		`,
		r.FirstName, r.LastName, r.PictureKey, r.UpdatedAt, r.ID,
	)
	if err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Exists(ctx context.Context, email string) (bool, error) {
	var exists bool
	if err := repo.pg.QueryRow(
		ctx,
		`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
		`,
		email,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("check if user email exists: %w", err)
	}

	return exists, nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*customer.Customer, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT 
			u.id, u.email, u.password_hash, u.verified_at, 
			c.first_name, c.last_name, c.picture_key, c.created_at, c.updated_at 
		FROM users u 
		INNER JOIN customers c ON c.id = u.id
		WHERE u.id = $1
		`,
		id,
	)

	var r customerRow
	if err := row.Scan(
		&r.ID,
		&r.Email,
		&r.PasswordHash,
		&r.VerifiedAt,
		&r.FirstName,
		&r.LastName,
		&r.PictureKey,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customer.ErrNotFound
		}

		return nil, fmt.Errorf("scan customer: %w", err)
	}

	return toEntity(r), nil
}

func (repo *PostgresRepo) FindByEmail(
	ctx context.Context,
	email string,
) (*customer.Customer, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT 
			u.id, u.email, u.password_hash, u.verified_at, 
			c.first_name, c.last_name, c.picture_key, c.created_at, c.updated_at 
		FROM users u 
		INNER JOIN customers c ON c.id = u.id 
		WHERE u.email = $1
		`,
		email,
	)

	var r customerRow
	if err := row.Scan(
		&r.ID,
		&r.Email,
		&r.PasswordHash,
		&r.VerifiedAt,
		&r.FirstName,
		&r.LastName,
		&r.PictureKey,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customer.ErrNotFound
		}

		return nil, fmt.Errorf("scan customer: %w", err)
	}

	return toEntity(r), nil
}
