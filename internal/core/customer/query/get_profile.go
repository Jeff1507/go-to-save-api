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

type ProfileOutput struct {
	ID         string
	Email      string
	FirstName  string
	LastName   string
	AvatarURL  *string
	VerifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type GetProfile struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewGetProfile(pg *pgxpool.Pool, storage storage.Storage) *GetProfile {
	return &GetProfile{
		pg:      pg,
		storage: storage,
	}
}

func (qry *GetProfile) Execute(ctx context.Context, id string) (ProfileOutput, error) {
	row := qry.pg.QueryRow(
		ctx,
		`
		SELECT 
			u.id, 
			u.email, 
			u.verified_at,
			c.first_name,
			c.last_name,
			c.avatar_key,
			c.created_at,
			c.updated_at
		FROM users u
		INNER JOIN customers c ON c.id = u.id
		WHERE u.id = $1
		`,
		id,
	)

	var (
		profile    ProfileOutput
		avatarKey  sql.NullString
		verifiedAt sql.NullTime
	)
	if err := row.Scan(
		&profile.ID,
		&profile.Email,
		&verifiedAt,
		&profile.FirstName,
		&profile.LastName,
		&avatarKey,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProfileOutput{}, customer.ErrNotFound
		}

		return ProfileOutput{}, fmt.Errorf("scan customer: %w", err)
	}

	if avatarKey.Valid {
		avatarURL := qry.storage.URL(avatarKey.String)
		profile.AvatarURL = &avatarURL
	}

	if verifiedAt.Valid {
		verifiedAt := verifiedAt.Time // Shadowing
		profile.VerifiedAt = &verifiedAt
	}

	return profile, nil
}
