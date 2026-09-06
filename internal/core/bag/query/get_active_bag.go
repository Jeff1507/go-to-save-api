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

type ActiveBagOutput struct {
	ID            string
	StoreID       string
	Title         string
	Description   *string
	Stock         int32
	Discount      int32
	OriginalPrice int64 // Preço em centavos
	FinalPrice    int64 // Preço em centavos
	PickupStart   time.Time
	PickupEnd     time.Time
	Items         []ItemOutput
}

type GetActiveBag struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewGetActiveBag(pg *pgxpool.Pool, storage storage.Storage) *GetActiveBag {
	return &GetActiveBag{
		pg:      pg,
		storage: storage,
	}
}

func (qry *GetActiveBag) Execute(ctx context.Context, id string) (ActiveBagOutput, error) {
	row := qry.pg.QueryRow(
		ctx,
		`
		SELECT
			id, 
			store_id,
			title, 
			description, 
			stock, 
			original_price, 
			final_price, 
			pickup_start, 
			pickup_end
		FROM bags
		WHERE id = $1 AND status = 'ACTIVE'
		`,
		id,
	)

	var (
		active      ActiveBagOutput
		description sql.NullString
	)
	if err := row.Scan(
		&active.ID,
		&active.StoreID,
		&active.Title,
		&description,
		&active.Stock,
		&active.OriginalPrice,
		&active.FinalPrice,
		&active.PickupStart,
		&active.PickupEnd,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ActiveBagOutput{}, bag.ErrNotFound
		}

		return ActiveBagOutput{}, fmt.Errorf("scan bag: %w", err)
	}

	if description.Valid {
		active.Description = &description.String
	}

	active.Discount = calculateDiscount(active.OriginalPrice, active.FinalPrice)

	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT id, name, image_key, quantity
		FROM bag_items
		WHERE bag_id = $1
		`,
		id,
	)
	if err != nil {
		return ActiveBagOutput{}, fmt.Errorf("query bag items: %w", err)
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
			return ActiveBagOutput{}, fmt.Errorf("scan bag item: %w", err)
		}

		if imageKey.Valid {
			imageURL := qry.storage.URL(imageKey.String)
			item.ImageURL = &imageURL
		}

		active.Items = append(active.Items, item)
	}

	if err := rows.Err(); err != nil {
		return ActiveBagOutput{}, fmt.Errorf("iterate rows: %w", err)
	}

	return active, nil
}
