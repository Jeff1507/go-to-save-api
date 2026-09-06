package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	Status        string
	PickupStart   time.Time
	PickupEnd     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ListBags struct {
	pg *pgxpool.Pool
}

func NewListBags(pg *pgxpool.Pool) *ListBags {
	return &ListBags{
		pg: pg,
	}
}

func (qry *ListBags) Execute(ctx context.Context, storeID string) ([]BagSummaryOutput, error) {
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
			status,
			pickup_start, 
			pickup_end, 
			created_at, 
			updated_at
		FROM bags
		WHERE store_id = $1
		ORDER BY created_at DESC
		`,
		storeID,
	)
	if err != nil {
		return nil, fmt.Errorf("query bags: %w", err)
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
			&b.Status,
			&b.PickupStart,
			&b.PickupEnd,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan bag: %w", err)
		}

		if description.Valid {
			b.Description = &description.String
		}

		b.Discount = calculateDiscount(b.OriginalPrice, b.FinalPrice)

		bags = append(bags, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bags: %w", err)
	}

	return bags, nil
}
