package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BagSummaryOutput struct {
	ID            string
	Title         string
	Description   *string
	Stock         int32
	Discount      int32
	OriginalPrice int64 // Preço em centavos
	FinalPrice    int64 // Preço em centavos
	PickupStart   time.Time
	PickupEnd     time.Time
}

type StoreOutput struct {
	ID          string
	Name        string
	CNPJ        string
	AvatarURL   *string
	Description *string
	CreatedAt   time.Time
	Rating      RatingOutput
	Address     AddressOutput
	Bags        []BagSummaryOutput
}

type GetStoreWithActiveBags struct {
	pg      *pgxpool.Pool
	storage storage.Storage
	clock   clock.Clock
}

func NewGetStoreWithActiveBags(
	pg *pgxpool.Pool,
	storage storage.Storage,
	clock clock.Clock,
) *GetStoreWithActiveBags {
	return &GetStoreWithActiveBags{
		pg:      pg,
		storage: storage,
		clock:   clock,
	}
}

func (qry *GetStoreWithActiveBags) Execute(ctx context.Context, id string) (StoreOutput, error) {
	now := qry.clock.Now()

	row := qry.pg.QueryRow(
		ctx,
		`
		SELECT
			u.id,
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
		FROM users u
		INNER JOIN stores s ON s.id = u.id
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE u.id = $1
		`,
		id,
	)

	var (
		s                 StoreOutput
		avatarKey         sql.NullString
		description       sql.NullString
		ratingSum         int32
		addressComplement sql.NullString
	)
	if err := row.Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return StoreOutput{}, store.ErrNotFound
		}

		return StoreOutput{}, fmt.Errorf("scan store: %w", err)
	}

	if avatarKey.Valid {
		avatarURL := qry.storage.URL(avatarKey.String)
		s.AvatarURL = &avatarURL
	}

	if description.Valid {
		s.Description = &description.String
	}

	s.Rating.Average = calculateAverageRating(ratingSum, s.Rating.Count)

	if addressComplement.Valid {
		s.Address.Complement = &addressComplement.String
	}

	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT
			id, 
			title, 
			description, 
			stock, 
			original_price, 
			final_price, 
			pickup_start, 
			pickup_end
		FROM bags
		WHERE store_id = $1 AND pickup_end > $2 AND status = 'ACTIVE' AND stock > 0
		`,
		id, now,
	)
	if err != nil {
		return StoreOutput{}, fmt.Errorf("query bags: %w", err)
	}
	defer rows.Close()

	bags := make([]BagSummaryOutput, 0)
	for rows.Next() {
		var (
			b           BagSummaryOutput
			description sql.NullString
		)
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&description,
			&b.Stock,
			&b.OriginalPrice,
			&b.FinalPrice,
			&b.PickupStart,
			&b.PickupEnd,
		); err != nil {
			return StoreOutput{}, fmt.Errorf("scan bag: %w", err)
		}

		if description.Valid {
			b.Description = &description.String
		}

		b.Discount = calculateDiscount(b.OriginalPrice, b.FinalPrice)

		bags = append(bags, b)
	}

	if err := rows.Err(); err != nil {
		return StoreOutput{}, fmt.Errorf("iterate rows: %w", err)
	}

	s.Bags = bags

	return s, nil
}

func calculateDiscount(originalPrice, finalPrice int64) int32 {
	if originalPrice <= 0 {
		return 0
	}

	return int32(100 - (finalPrice * 100 / originalPrice))
}
