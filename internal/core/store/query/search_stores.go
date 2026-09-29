package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchStoresInput struct {
	Name  string
	Limit int32
}

type StoreResultOutput struct {
	ID          string
	Name        string
	CNPJ        string
	AvatarURL   *string
	Description *string
	CreatedAt   time.Time
	Rating      RatingOutput
	Address     AddressOutput
}

type SearchStores struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewSearchStores(pg *pgxpool.Pool, storage storage.Storage) *SearchStores {
	return &SearchStores{
		pg:      pg,
		storage: storage,
	}
}

func (qry *SearchStores) Execute(
	ctx context.Context,
	input SearchStoresInput,
) ([]StoreResultOutput, error) {
	limit := input.Limit
	if limit <= 0 || limit > 5 {
		limit = 5
	}

	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT
			s.id,
			s.name,
			s.cnpj,
			s.avatar_key,
			s.description,
			s.created_at,
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
		FROM stores s
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE unaccent(s.name) ILIKE unaccent('%' || $1 || '%')
		ORDER BY s.name
		LIMIT $2
		`,
		input.Name, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query stores: %w", err)
	}
	defer rows.Close()

	stores := make([]StoreResultOutput, 0)
	for rows.Next() {
		var (
			s                 StoreResultOutput
			avatarKey         sql.NullString
			description       sql.NullString
			ratingSum         int32
			addressComplement sql.NullString
		)
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.CNPJ,
			&avatarKey,
			&description,
			&s.CreatedAt,
			&ratingSum,
			&s.Rating.Count,
			&s.Address.PostalCode,
			&s.Address.Street,
			&s.Address.Number,
			&s.Address.Neighborhood,
			&addressComplement,
			&s.Address.City,
			&s.Address.State,
			&s.Address.Latitude,
			&s.Address.Longitude,
		); err != nil {
			return nil, fmt.Errorf("scan store: %w", err)
		}

		if avatarKey.Valid {
			avatarURL := qry.storage.URL(avatarKey.String)
			s.AvatarURL = &avatarURL
		}

		if description.Valid {
			s.Description = &description.String
		}

		if addressComplement.Valid {
			s.Address.Complement = &addressComplement.String
		}

		s.Rating.Average = calculateAverageRating(ratingSum, s.Rating.Count)

		stores = append(stores, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return stores, nil
}
