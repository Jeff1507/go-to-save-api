package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderSummaryOutput struct {
	ID          string
	Status      string
	TotalPrice  int64 // Preço em centavos
	PickupStart time.Time
	PickupEnd   time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
	Store       StoreOutput
}

type ListOrders struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewListOrders(pg *pgxpool.Pool, storage storage.Storage) *ListOrders {
	return &ListOrders{
		pg:      pg,
		storage: storage,
	}
}

func (qry *ListOrders) Execute(
	ctx context.Context,
	customerID string,
) ([]OrderSummaryOutput, error) {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT
			o.id,
			o.status,
			o.total_price,
			o.pickup_start,
			o.pickup_end,
			o.created_at,
			o.completed_at,
			s.id, 
			s.name, 
			s.avatar_key,
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
		FROM orders o
		INNER JOIN stores s ON s.id = o.store_id
		INNER JOIN stores_rating sr ON sr.store_id = s.id
		INNER JOIN addresses a ON a.store_id = s.id
		WHERE o.customer_id = $1
		ORDER BY o.created_at DESC
		`,
		customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	orders := make([]OrderSummaryOutput, 0)
	for rows.Next() {
		var (
			ord                    OrderSummaryOutput
			completedAt            sql.NullTime
			storeAvatarKey         sql.NullString
			storeRatingSum         int32
			storeAddressComplement sql.NullString
		)
		if err := rows.Scan(
			&ord.ID,
			&ord.Status,
			&ord.TotalPrice,
			&ord.PickupStart,
			&ord.PickupEnd,
			&ord.CreatedAt,
			&ord.CompletedAt,
			&ord.Store.ID,
			&ord.Store.Name,
			&storeAvatarKey,
			&storeRatingSum,
			&ord.Store.Rating.Count,
			&ord.Store.Address.PostalCode,
			&ord.Store.Address.Street,
			&ord.Store.Address.Number,
			&ord.Store.Address.Neighborhood,
			&storeAddressComplement,
			&ord.Store.Address.City,
			&ord.Store.Address.State,
			&ord.Store.Address.Latitude,
			&ord.Store.Address.Longitude,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		if completedAt.Valid {
			ord.CompletedAt = &completedAt.Time
		}

		if storeAvatarKey.Valid {
			imageURL := qry.storage.URL(storeAvatarKey.String)
			ord.Store.AvatarURL = &imageURL
		}

		ord.Store.Rating.Average = calculateAverageRating(storeRatingSum, ord.Store.Rating.Count)

		if storeAddressComplement.Valid {
			ord.Store.Address.Complement = &storeAddressComplement.String
		}

		orders = append(orders, ord)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return orders, nil
}
