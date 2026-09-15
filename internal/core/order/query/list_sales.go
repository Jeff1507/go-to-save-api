package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListSalesInput struct {
	StoreID string
	Status  string
}

type SaleSummaryOutput struct {
	ID          string
	Status      string
	TotalPrice  int64
	PickupStart time.Time
	PickupEnd   time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
	Customer    CustomerOutput
}

// Sale = Order
type ListSales struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewListSales(pg *pgxpool.Pool, storage storage.Storage) *ListSales {
	return &ListSales{
		pg:      pg,
		storage: storage,
	}
}

func (qry *ListSales) Execute(
	ctx context.Context,
	input ListSalesInput,
) ([]SaleSummaryOutput, error) {
	if input.Status != "" {
		if _, err := order.ParseStatus(input.Status); err != nil {
			return nil, fmt.Errorf("invalid status: %w", err)
		}
	}

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
			c.id,
			c.first_name,
			c.last_name,
			c.avatar_key
		FROM orders o
		INNER JOIN customers c ON c.id = o.customer_id
		WHERE o.store_id = $1 AND ($2::text = '' OR o.status = $2::order_status)
		ORDER BY o.created_at DESC
		`,
		input.StoreID, input.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	sales := make([]SaleSummaryOutput, 0)
	for rows.Next() {
		var (
			sale              SaleSummaryOutput
			completedAt       sql.NullTime
			customerFirstName string
			customerLastName  string
			customerAvatarKey sql.NullString
		)
		if err := rows.Scan(
			&sale.ID,
			&sale.Status,
			&sale.TotalPrice,
			&sale.PickupStart,
			&sale.PickupEnd,
			&sale.CreatedAt,
			&completedAt,
			&sale.Customer.ID,
			&customerFirstName,
			&customerLastName,
			&customerAvatarKey,
		); err != nil {
			return nil, fmt.Errorf("scan sale: %w", err)
		}

		if completedAt.Valid {
			sale.CompletedAt = &completedAt.Time
		}

		if customerAvatarKey.Valid {
			avatarURL := qry.storage.URL(customerAvatarKey.String)
			sale.Customer.AvatarURL = &avatarURL
		}

		sale.Customer.FullName = customerFirstName + " " + customerLastName

		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return sales, nil
}
