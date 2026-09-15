package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileOutput struct {
	ID          string
	Email       string
	Name        string
	CNPJ        string
	AvatarURL   *string
	Description *string
	VerifiedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Rating      RatingOutput
	Address     AddressOutput
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

	var (
		profile           ProfileOutput
		avatarKey         sql.NullString
		description       sql.NullString
		verifiedAt        sql.NullTime
		ratingSum         int32
		addressComplement sql.NullString
	)
	if err := row.Scan(
		&profile.ID,
		&profile.Email,
		&verifiedAt,
		&profile.Name,
		&profile.CNPJ,
		&avatarKey,
		&description,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&ratingSum,
		&profile.Rating.Count,
		&profile.Address.PostalCode,
		&profile.Address.Street,
		&profile.Address.Number,
		&profile.Address.Neighborhood,
		&addressComplement,
		&profile.Address.City,
		&profile.Address.State,
		&profile.Address.Latitude,
		&profile.Address.Longitude,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProfileOutput{}, store.ErrNotFound
		}

		return ProfileOutput{}, fmt.Errorf("scan store: %w", err)
	}

	if avatarKey.Valid {
		avatarURL := qry.storage.URL(avatarKey.String)
		profile.AvatarURL = &avatarURL
	}

	if description.Valid {
		profile.Description = &description.String
	}

	if verifiedAt.Valid {
		profile.VerifiedAt = &verifiedAt.Time
	}

	if addressComplement.Valid {
		profile.Address.Complement = &addressComplement.String
	}

	profile.Rating.Average = calculateAverageRating(ratingSum, profile.Rating.Count)

	return profile, nil
}
