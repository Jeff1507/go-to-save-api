package repository

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
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

func (repo *PostgresRepo) Save(ctx context.Context, s *store.Store) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	r := toRow(s)

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
			return store.ErrEmailAlreadyInUse
		}

		return fmt.Errorf("insert user: %v", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO stores(id, name, cnpj, picture_key, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		r.ID, r.Name, r.CNPJ, r.PictureKey, r.Description, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return store.ErrDuplicateCNPJ
		}

		return fmt.Errorf("insert store: %v", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO addresses(
			store_id, postal_code, street, number, neighborhood, complement, city, state, 
			latitude, longitude
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`,
		r.ID, r.Address.PostalCode, r.Address.Street, r.Address.Number, r.Address.Neighborhood,
		r.Address.Complement, r.Address.City, r.Address.State,
		r.Address.Latitude, r.Address.Longitude,
	)
	if err != nil {
		return fmt.Errorf("insert address: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %v", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, s *store.Store) error {
	return nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*store.Store, error) {
	return nil, nil
}

func (repo *PostgresRepo) FindByEmail(ctx context.Context, email string) (*store.Store, error) {
	return nil, nil
}
