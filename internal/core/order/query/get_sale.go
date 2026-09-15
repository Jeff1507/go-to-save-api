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

type GetSaleInput struct {
	ID      string
	StoreID string
}

type SaleOutput struct {
	ID          string
	Status      string
	TotalPrice  int64
	PickupStart time.Time
	PickupEnd   time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
	Customer    CustomerOutput
	Bags        []BagSummaryOutput
}

type GetSale struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewGetSale(pg *pgxpool.Pool, storage storage.Storage) *GetSale {
	return &GetSale{
		pg:      pg,
		storage: storage,
	}
}

func (qry *GetSale) Execute(ctx context.Context, input GetSaleInput) (SaleOutput, error) {
	sale, err := qry.fetchSale(ctx, input.ID, input.StoreID)
	if err != nil {
		return SaleOutput{}, fmt.Errorf("fetch order: %w", err)
	}

	bags, err := qry.fetchBags(ctx, input.ID)
	if err != nil {
		return SaleOutput{}, fmt.Errorf("fetch bags: %w", err)
	}

	bagIDs := make([]string, len(bags))
	bagsByID := make(map[string]*BagSummaryOutput, len(bags))
	for i := range bags {
		bagIDs[i] = bags[i].ID
		bagsByID[bags[i].ID] = &bags[i]
	}

	if err := qry.fillBagItems(ctx, bagIDs, bagsByID); err != nil {
		return SaleOutput{}, fmt.Errorf("fill bag items: %w", err)
	}

	sale.Bags = bags

	return sale, nil
}

func (qry *GetSale) fetchSale(ctx context.Context, id, storeID string) (SaleOutput, error) {
	row := qry.pg.QueryRow(
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
			c.id,
			c.first_name,
			c.last_name,
			c.avatar_key
		FROM orders o
		INNER JOIN customers c ON c.id = o.customer_id
		WHERE o.id = $1 AND o.store_id = $2 
		`,
		id, storeID,
	)

	var (
		sale        SaleOutput
		completedAt sql.NullTime
		cust        struct {
			FirstName string
			LastName  string
			AvatarKey sql.NullString
		}
	)
	if err := row.Scan(
		&sale.ID,
		&sale.Status,
		&sale.TotalPrice,
		&sale.PickupStart,
		&sale.PickupEnd,
		&sale.CreatedAt,
		&completedAt,
		&sale.Customer.ID,
		&cust.FirstName,
		&cust.LastName,
		&cust.AvatarKey,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SaleOutput{}, order.ErrNotFound
		}

		return SaleOutput{}, fmt.Errorf("scan sale: %w", err)
	}

	if completedAt.Valid {
		sale.CompletedAt = &completedAt.Time
	}

	if cust.AvatarKey.Valid {
		avatarURL := qry.storage.URL(cust.AvatarKey.String)
		sale.Customer.AvatarURL = &avatarURL
	}

	sale.Customer.FullName = cust.FirstName + " " + cust.LastName

	return sale, nil
}

func (qry *GetSale) fetchBags(ctx context.Context, id string) ([]BagSummaryOutput, error) {
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

func (qry *GetSale) fillBagItems(
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
