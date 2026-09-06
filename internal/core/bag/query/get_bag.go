package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GetBagInput struct {
	ID      string
	StoreID string
}

type BagOutput struct {
	ID            string
	Title         string
	Description   *string
	Stock         int32
	Discount      int32
	OriginalPrice int64
	FinalPrice    int64
	Status        string
	PickupStart   time.Time
	PickupEnd     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Items         []ItemOutput
}

type GetBag struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewGetBag(pg *pgxpool.Pool, storage storage.Storage) *GetBag {
	return &GetBag{
		pg:      pg,
		storage: storage,
	}
}

func (qry *GetBag) Execute(ctx context.Context, input GetBagInput) (BagOutput, error) {
	row := qry.pg.QueryRow(
		ctx,
		`
		SELECT
			id, 
			title, 
			description, 
			stock, 
			original_price, 
			final_price, 
			status,
			pickup_start, 
			pickup_end, 
			created_at, 
			updated_at
		FROM bags
		WHERE id = $1 AND store_id = $2
		`,
		input.ID, input.StoreID,
	)

	var (
		b           BagOutput
		description sql.NullString
	)
	if err := row.Scan(
		&b.ID,
		&b.Title,
		&description,
		&b.Stock,
		&b.OriginalPrice,
		&b.FinalPrice,
		&b.Status,
		&b.PickupStart,
		&b.PickupEnd,
		&b.CreatedAt,
		&b.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BagOutput{}, bag.ErrNotFound
		}

		return BagOutput{}, fmt.Errorf("scan bag: %w", err)
	}

	if description.Valid {
		b.Description = &description.String
	}

	b.Discount = calculateDiscount(b.OriginalPrice, b.FinalPrice)

	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT id, name, image_key, quantity
		FROM bag_items
		WHERE bag_id = $1
		`,
		input.ID,
	)
	if err != nil {
		return BagOutput{}, fmt.Errorf("query bag items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			item     ItemOutput
			imageKey sql.NullString
		)
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&imageKey,
			&item.Quantity,
		); err != nil {
			return BagOutput{}, fmt.Errorf("scan bag item: %w", err)
		}

		if imageKey.Valid {
			imageURL := qry.storage.URL(imageKey.String)
			item.ImageURL = &imageURL
		}

		b.Items = append(b.Items, item)
	}

	if err := rows.Err(); err != nil {
		return BagOutput{}, fmt.Errorf("iterate rows: %w", err)
	}

	return b, nil
}
