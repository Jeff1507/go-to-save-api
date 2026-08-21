package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductOutput struct {
	ID        string
	Name      string
	ImageURL  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListProducts struct {
	pg      *pgxpool.Pool
	storage storage.Storage
}

func NewListProducts(pg *pgxpool.Pool, storage storage.Storage) *ListProducts {
	return &ListProducts{
		pg:      pg,
		storage: storage,
	}
}

func (qry *ListProducts) Execute(ctx context.Context, storeID string) ([]ProductOutput, error) {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT id, name, image_key, created_at, updated_at
		FROM products
		WHERE store_id = $1
		ORDER BY name ASC
		`,
		storeID,
	)
	if err != nil {
		return nil, fmt.Errorf("select products: %w", err)
	}
	defer rows.Close()

	products := make([]ProductOutput, 0)
	for rows.Next() {
		var (
			prod     ProductOutput
			imageKey sql.NullString
		)
		if err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&imageKey,
			&prod.CreatedAt,
			&prod.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		if imageKey.Valid {
			imageURL := qry.storage.URL(imageKey.String)
			prod.ImageURL = &imageURL
		}

		products = append(products, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return products, nil
}
