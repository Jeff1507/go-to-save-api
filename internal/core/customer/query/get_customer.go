package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerProfile struct {
	ID         string
	Email      string
	FirstName  string
	LastName   string
	PictureURL *string
	VerifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type GetCustomer struct {
	pg          *pgxpool.Pool
	fileStorage storage.Storage
}

func NewGetCustomer(pg *pgxpool.Pool, fileStorage storage.Storage) *GetCustomer {
	return &GetCustomer{
		pg:          pg,
		fileStorage: fileStorage,
	}
}

func (qry *GetCustomer) Execute(ctx context.Context, id string) (CustomerProfile, error) {
	row := qry.pg.QueryRow(
		ctx,
		`
		SELECT 
			u.id, 
			u.email, 
			u.verified_at,
			c.first_name,
			c.last_name,
			c.picture_key,
			c.created_at,
			c.updated_at
		FROM users u
		INNER JOIN customers c ON c.id = u.id
		WHERE u.id = $1
		`,
		id,
	)

	var (
		profile    CustomerProfile
		pictureKey sql.NullString
		verifiedAt sql.NullTime
	)
	if err := row.Scan(
		&profile.ID,
		&profile.Email,
		&verifiedAt,
		&profile.FirstName,
		&profile.LastName,
		&pictureKey,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CustomerProfile{}, customer.ErrNotFound
		}

		return CustomerProfile{}, fmt.Errorf("scan customer: %w", err)
	}

	if pictureKey.Valid {
		pictureURL := qry.fileStorage.URL(pictureKey.String)
		profile.PictureURL = &pictureURL
	}

	if verifiedAt.Valid {
		verifiedAt := verifiedAt.Time // Shadowing
		profile.VerifiedAt = &verifiedAt
	}

	return profile, nil
}
