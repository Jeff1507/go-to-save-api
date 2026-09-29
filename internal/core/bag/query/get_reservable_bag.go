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

type ReservableBagOutput struct {
	ID            string
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

type GetReservableBag struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewGetReservableBag(pg *pgxpool.Pool, storage storage.Storage) *GetReservableBag {
	return &GetReservableBag{
		pg:      pg,
		storage: storage,
	}
}

func (qry *GetReservableBag) Execute(ctx context.Context, id string) (ReservableBagOutput, error) {
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
			pickup_start, 
			pickup_end
		FROM bags
		WHERE id = $1
		`,
		id,
	)

	var (
		b           ReservableBagOutput
		description sql.NullString
	)
	if err := row.Scan(
		&b.ID,
		&b.Title,
		&description,
		&b.Stock,
		&b.OriginalPrice,
		&b.FinalPrice,
		&b.PickupStart,
		&b.PickupEnd,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReservableBagOutput{}, bag.ErrNotFound
		}

		return ReservableBagOutput{}, fmt.Errorf("scan bag: %w", err)
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
		id,
	)
	if err != nil {
		return ReservableBagOutput{}, fmt.Errorf("query bag items: %w", err)
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
			return ReservableBagOutput{}, fmt.Errorf("scan bag item: %w", err)
		}

		if imageKey.Valid {
			imageURL := qry.storage.URL(imageKey.String)
			item.ImageURL = &imageURL
		}

		b.Items = append(b.Items, item)
	}

	if err := rows.Err(); err != nil {
		return ReservableBagOutput{}, fmt.Errorf("iterate rows: %w", err)
	}

	return b, nil
}
