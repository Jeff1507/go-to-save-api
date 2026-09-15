package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GetOrderInput struct {
	ID         string
	CustomerID string
}

type ReviewOutput struct {
	ID         string
	Rating     int16
	Comment    *string
	ReviewedAt time.Time
}

type OrderOutput struct {
	ID               string
	Status           string
	TotalPrice       int64
	ConfirmationCode string
	CreatedAt        time.Time
	CompletedAt      *time.Time
	PickupStart      time.Time
	PickupEnd        time.Time
	Store            StoreOutput
	Bags             []BagSummaryOutput
	Review           *ReviewOutput
}

type GetOrder struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewGetOrder(pg *pgxpool.Pool, storage storage.Storage) *GetOrder {
	return &GetOrder{
		pg:      pg,
		storage: storage,
	}
}
func (qry *GetOrder) Execute(ctx context.Context, input GetOrderInput) (OrderOutput, error) {
	ord, err := qry.fetchOrder(ctx, input.ID, input.CustomerID)
	if err != nil {
		return OrderOutput{}, fmt.Errorf("fetch order: %w", err)
	}

	bags, err := qry.fetchBags(ctx, input.ID)
	if err != nil {
		return OrderOutput{}, fmt.Errorf("fetch bags: %w", err)
	}

	bagIDs := make([]string, len(bags))
	bagsByID := make(map[string]*BagSummaryOutput, len(bags))
	for i := range bags {
		bagIDs[i] = bags[i].ID
		bagsByID[bags[i].ID] = &bags[i]
	}

	if err := qry.fillBagItems(ctx, bagIDs, bagsByID); err != nil {
		return OrderOutput{}, fmt.Errorf("fill bag items: %w", err)
	}

	ord.Bags = bags

	return ord, nil
}

func (qry *GetOrder) fetchOrder(ctx context.Context, id, customerID string) (OrderOutput, error) {
	row := qry.pg.QueryRow(
		ctx,
		`
		SELECT
			o.id,
			o.status,
			o.total_price,
			o.confirmation_code,
			o.pickup_start,
			o.pickup_end,
			o.created_at,
			o.completed_at,
			s.id,
			s.name,
			s.avatar_key,
			a.postal_code,
			a.street,
			a.number,
			a.neighborhood,
			a.complement,
			a.city,
			a.state,
			a.latitude,
			a.longitude,
			r.id,
			r.rating,
			r.comment,
			r.reviewed_at
		FROM orders o
		INNER JOIN stores s ON s.id = o.store_id
		INNER JOIN addresses a ON a.store_id = s.id
		LEFT JOIN reviews r ON r.order_id = o.id
		WHERE o.id = $1 AND o.customer_id = $2 
		`,
		id, customerID,
	)

	var (
		ord                    OrderOutput
		completedAt            sql.NullTime
		storeAvatarKey         sql.NullString
		storeAddressComplement sql.NullString
		review                 struct {
			ID         sql.NullString
			Rating     sql.NullInt16
			Comment    sql.NullString
			ReviewedAt sql.NullTime
		}
	)
	if err := row.Scan(
		&ord.ID,
		&ord.Status,
		&ord.TotalPrice,
		&ord.ConfirmationCode,
		&ord.PickupStart,
		&ord.PickupEnd,
		&ord.CreatedAt,
		&completedAt,
		&ord.Store.ID,
		&ord.Store.Name,
		&storeAvatarKey,
		&ord.Store.Address.PostalCode,
		&ord.Store.Address.Street,
		&ord.Store.Address.Number,
		&ord.Store.Address.Neighborhood,
		&storeAddressComplement,
		&ord.Store.Address.City,
		&ord.Store.Address.State,
		&ord.Store.Address.Latitude,
		&ord.Store.Address.Longitude,
		&review.ID,
		&review.Rating,
		&review.Comment,
		&review.ReviewedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderOutput{}, order.ErrNotFound
		}

		return OrderOutput{}, fmt.Errorf("scan order: %w", err)
	}

	if completedAt.Valid {
		ord.CompletedAt = &completedAt.Time
	}

	if storeAvatarKey.Valid {
		avatarURL := qry.storage.URL(storeAvatarKey.String)
		ord.Store.AvatarURL = &avatarURL
	}

	if storeAddressComplement.Valid {
		ord.Store.Address.Complement = &storeAddressComplement.String
	}

	if review.ID.Valid {
		ord.Review = &ReviewOutput{
			ID:         review.ID.String,
			Rating:     review.Rating.Int16,
			Comment:    &review.Comment.String,
			ReviewedAt: review.ReviewedAt.Time,
		}
	}

	return ord, nil
}

func (qry *GetOrder) fetchBags(ctx context.Context, id string) ([]BagSummaryOutput, error) {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT oi.bag_id, b.title, oi.quantity, oi.unit_price
		FROM order_items oi
		INNER JOIN bags b ON b.id = oi.bag_id
		WHERE oi.order_id = $1
		`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	bags := make([]BagSummaryOutput, 0)
	for rows.Next() {
		var b BagSummaryOutput
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Quantity,
			&b.UnitPrice,
		); err != nil {
			return nil, fmt.Errorf("scan bag: %w", err)
		}

		b.Items = make([]BagItemOutput, 0)
		bags = append(bags, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return bags, nil
}

func (qry *GetOrder) fillBagItems(
	ctx context.Context,
	bagIDs []string,
	bagsByID map[string]*BagSummaryOutput,
) error {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT id, bag_id, name, image_key, quantity
		FROM bag_items
		WHERE bag_id = ANY($1)
		`,
		bagIDs,
	)
	if err != nil {
		return fmt.Errorf("query bag items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			item     BagItemOutput
			bagID    string
			imageKey sql.NullString
		)
		if err := rows.Scan(
			&item.ID,
			&bagID,
			&item.Name,
			&imageKey,
			&item.Quantity,
		); err != nil {
			return fmt.Errorf("scan bag item: %w", err)
		}

		if imageKey.Valid {
			imageURL := qry.storage.URL(imageKey.String)
			item.ImageURL = &imageURL
		}

		if b, ok := bagsByID[bagID]; ok {
			b.Items = append(b.Items, item)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate rows: %w", err)
	}

	return nil
}
