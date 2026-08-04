package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreProduct struct {
	ID        string
	Name      string
	ImageURL  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListStoreProducts struct {
	pg          *pgxpool.Pool
	fileStorage storage.Storage
}

func NewListStoreProducts(pg *pgxpool.Pool, fileStorage storage.Storage) *ListStoreProducts {
	return &ListStoreProducts{
		pg:          pg,
		fileStorage: fileStorage,
	}
}

func (qry *ListStoreProducts) Execute(
	ctx context.Context,
	storeID string,
) ([]StoreProduct, error) {
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
		return nil, fmt.Errorf("select products: %v", err)
	}
	defer rows.Close()

	products := make([]StoreProduct, 0)
	for rows.Next() {
		var (
			p      StoreProduct
			imgKey sql.NullString
		)
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&imgKey,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product: %v", err)
		}

		if imgKey.Valid {
			imgURL := qry.fileStorage.URL(imgKey.String)
			p.ImageURL = &imgURL
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %v", err)
	}

	return products, nil
}
