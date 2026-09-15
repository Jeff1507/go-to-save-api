package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
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

func (repo *PostgresRepo) Save(ctx context.Context, s *store.Store) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
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

		return fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO stores(id, name, cnpj, avatar_key, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		r.ID, r.Name, r.CNPJ, r.AvatarKey, r.Description, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return store.ErrDuplicateCNPJ
		}

		return fmt.Errorf("insert store: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO stores_rating(store_id, sum, count)
		VALUES ($1, $2, $3)
		`,
		r.ID, r.RatingSum, r.RatingCount,
	)
	if err != nil {
		return fmt.Errorf("insert store rating: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO addresses(
			store_id, 
			postal_code, 
			street, 
			number, 
			neighborhood, 
			complement, 
			city, 
			state, 
			latitude, 
			longitude
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`,
		r.ID,
		r.Address.PostalCode,
		r.Address.Street,
		r.Address.Number,
		r.Address.Neighborhood,
		r.Address.Complement,
		r.Address.City,
		r.Address.State,
		r.Address.Latitude,
		r.Address.Longitude,
	)
	if err != nil {
		return fmt.Errorf("insert address: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, s *store.Store) error {
	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	r := toRow(s)

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
			return store.ErrEmailAlreadyInUse
		}

		return fmt.Errorf("update user: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE stores SET name = $1, cnpj = $2, avatar_key = $3, description = $4, updated_at = $5
		WHERE id = $6
		`,
		r.Name, r.CNPJ, r.AvatarKey, r.Description, r.UpdatedAt, r.ID,
	)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return store.ErrDuplicateCNPJ
		}

		return fmt.Errorf("update store: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE stores_rating
		SET sum = $1, count = $2
		WHERE store_id = $3
		`,
		r.RatingSum, r.RatingCount, r.ID,
	)
	if err != nil {
		return fmt.Errorf("update store rating: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE addresses 
		SET 
			postal_code = $1, 
			street = $2, 
			number = $3, 
			neighborhood = $4, 
			complement = $5, 
			city = $6, 
			state = $7, 
			latitude = $8, 
			longitude = $9
		WHERE store_id = $10
		`,
		r.Address.PostalCode,
		r.Address.Street,
		r.Address.Number,
		r.Address.Neighborhood,
		r.Address.Complement,
		r.Address.City,
		r.Address.State,
		r.Address.Latitude,
		r.Address.Longitude,
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("update address: %w", err)
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

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*store.Store, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT 
			u.id,
			u.email,
			u.password_hash,
			u.verified_at,
			s.name,
			s.cnpj,
			s.avatar_key,
			s.description,
			s.created_at,
			s.updated_at,
			sr.sum,
			sr.count,
			a.postal_code,
			a.street,
			a.number,
			a.neighborhood,
			a.complement,
			a.city,
			a.state,
			a.latitude,
			a.longitude
		FROM users u
		INNER JOIN stores s ON s.id = u.id
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE u.id = $1
		`,
		id,
	)

	var r storeRow
	if err := row.Scan(
		&r.ID,
		&r.Email,
		&r.PasswordHash,
		&r.VerifiedAt,
		&r.Name,
		&r.CNPJ,
		&r.AvatarKey,
		&r.Description,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.RatingSum,
		&r.RatingCount,
		&r.Address.PostalCode,
		&r.Address.Street,
		&r.Address.Number,
		&r.Address.Neighborhood,
		&r.Address.Complement,
		&r.Address.City,
		&r.Address.State,
		&r.Address.Latitude,
		&r.Address.Longitude,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, fmt.Errorf("scan store: %w", err)
	}

	return toEntity(r), nil
}

func (repo *PostgresRepo) FindByEmail(ctx context.Context, email string) (*store.Store, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT 
			u.id,
			u.email,
			u.password_hash,
			u.verified_at,
			s.name,
			s.cnpj,
			s.avatar_key,
			s.description,
			s.created_at,
			s.updated_at,
			sr.sum,
			sr.count,
			a.postal_code,
			a.street,
			a.number,
			a.neighborhood,
			a.complement,
			a.city,
			a.state,
			a.latitude,
			a.longitude
		FROM users u
		INNER JOIN stores s ON s.id = u.id
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE u.email = $1
		`,
		email,
	)

	var r storeRow
	if err := row.Scan(
		&r.ID,
		&r.Email,
		&r.PasswordHash,
		&r.VerifiedAt,
		&r.Name,
		&r.CNPJ,
		&r.AvatarKey,
		&r.Description,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.RatingSum,
		&r.RatingCount,
		&r.Address.PostalCode,
		&r.Address.Street,
		&r.Address.Number,
		&r.Address.Neighborhood,
		&r.Address.Complement,
		&r.Address.City,
		&r.Address.State,
		&r.Address.Latitude,
		&r.Address.Longitude,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, fmt.Errorf("scan store: %w", err)
	}

	return toEntity(r), nil
}

func (repo *PostgresRepo) FindAndUpdate(
	ctx context.Context,
	id string,
	fn func(s *store.Store) error,
) error {
	tx, ok := transaction.Querier(ctx)
	if ok {
		return repo.findAndUpdate(ctx, tx, id, fn)
	}

	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := repo.findAndUpdate(ctx, tx, id, fn); err != nil {
		return fmt.Errorf("find and update with transaction: %w", err)
	}

	return tx.Commit(ctx)
}

// Por enquanto, só atualiza o store_rating
// Se quiser expandir para mais campos -> atualizar a query de update
func (repo *PostgresRepo) findAndUpdate(
	ctx context.Context,
	tx pgx.Tx,
	id string,
	fn func(s *store.Store) error,
) error {
	row := tx.QueryRow(
		ctx,
		`
		SELECT 
			u.id,
			u.email,
			u.password_hash,
			u.verified_at,
			s.name,
			s.cnpj,
			s.avatar_key,
			s.description,
			s.created_at,
			s.updated_at,
			sr.sum,
			sr.count,
			a.postal_code,
			a.street,
			a.number,
			a.neighborhood,
			a.complement,
			a.city,
			a.state,
			a.latitude,
			a.longitude
		FROM users u
		INNER JOIN stores s ON s.id = u.id
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE u.id = $1 FOR UPDATE
		`,
		id,
	)

	var r storeRow
	if err := row.Scan(
		&r.ID,
		&r.Email,
		&r.PasswordHash,
		&r.VerifiedAt,
		&r.Name,
		&r.CNPJ,
		&r.AvatarKey,
		&r.Description,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.RatingSum,
		&r.RatingCount,
		&r.Address.PostalCode,
		&r.Address.Street,
		&r.Address.Number,
		&r.Address.Neighborhood,
		&r.Address.Complement,
		&r.Address.City,
		&r.Address.State,
		&r.Address.Latitude,
		&r.Address.Longitude,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return store.ErrNotFound
		}

		return fmt.Errorf("scan store: %w", err)
	}

	s := toEntity(r)

	if err := fn(s); err != nil {
		return fmt.Errorf("execute callback function: %w", err)
	}

	_, err := tx.Exec(
		ctx,
		`
		UPDATE stores_rating
		SET sum = $1, count = $2
		WHERE store_id = $3
		`,
		s.Rating().Sum(), s.Rating().Count(), s.ID(),
	)
	if err != nil {
		return fmt.Errorf("update store rating: %w", err)
	}

	return nil
}
