package query

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListNearbyStoresInput struct {
	Latitude  float64
	Longitude float64
	Radius    int32
}

type StoreSummaryOutput struct {
	ID           string
	Name         string
	CNPJ         string
	AvatarURL    *string
	Description  *string
	BagsCount    int32
	BagsMinPrice int64
	DistanceKM   float64
	CreatedAt    time.Time
	Rating       RatingOutput
	Address      AddressOutput
}

type ListNearbyStores struct {
	pg      *pgxpool.Pool
	storage storage.Storage
	clock   clock.Clock
}

func NewListNearbyStores(
	pg *pgxpool.Pool,
	storage storage.Storage,
	clock clock.Clock,
) *ListNearbyStores {
	return &ListNearbyStores{
		pg:      pg,
		storage: storage,
		clock:   clock,
	}
}

func (qry *ListNearbyStores) Execute(
	ctx context.Context,
	input ListNearbyStoresInput,
) ([]StoreSummaryOutput, error) {
	now := qry.clock.Now()

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
			a.longitude,
			(
				SELECT COUNT(*) FROM bags b 
				WHERE b.store_id = s.id 
				AND b.status = 'ACTIVE' 
				AND b.stock > 0 
				AND b.pickup_end > $4
			) AS bags_count,
			(
				SELECT MIN(b.final_price) FROM bags b 
				WHERE b.store_id = s.id 
				AND b.status = 'ACTIVE' 
				AND b.stock > 0 
				AND b.pickup_end > $4
			) AS bags_min_price,
			ST_Distance(
				a.location,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
			) AS distance
		FROM stores s
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE ST_DWithin(
			a.location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3
		)
		AND EXISTS (
			SELECT 1 FROM bags b 
			WHERE b.store_id = s.id 
			AND b.status = 'ACTIVE' 
			AND b.stock > 0 
			AND b.pickup_end > $4
		)
		ORDER BY distance
		`,
		input.Longitude, input.Latitude, input.Radius, now,
	)
	if err != nil {
		return nil, fmt.Errorf("query stores: %w", err)
	}
	defer rows.Close()

	stores := make([]StoreSummaryOutput, 0)
	for rows.Next() {
		var (
			s                 StoreSummaryOutput
			avatarKey         sql.NullString
			description       sql.NullString
			ratingSum         int32
			addressComplement sql.NullString
			distanceMeters    float64
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
			&s.BagsCount,
			&s.BagsMinPrice,
			&distanceMeters,
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

		s.DistanceKM = formatDistanceKM(distanceMeters)

		s.Rating.Average = calculateAverageRating(ratingSum, s.Rating.Count)

		stores = append(stores, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return stores, nil
}

func formatDistanceKM(meters float64) float64 {
	return math.Round(meters/1000*10) / 10
}
