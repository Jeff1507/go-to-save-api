package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreBagPickupWindow struct {
	Start time.Time
	End   time.Time
}

type StoreBagItem struct {
	ID       string
	Name     string
	ImageURL *string
	Quantity int32
}

type StoreBag struct {
	ID            string
	Title         string
	Description   *string
	Stock         int32
	OriginalPrice int64
	FinalPrice    int64
	Status        string
	PickupWindow  StoreBagPickupWindow
	Items         []StoreBagItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ListStoreBags struct {
	pg          *pgxpool.Pool
	fileStorage storage.Storage
}

func NewListStoreBags(pg *pgxpool.Pool, fileStorage storage.Storage) *ListStoreBags {
	return &ListStoreBags{
		pg:          pg,
		fileStorage: fileStorage,
	}
}

func (qry *ListStoreBags) Execute(ctx context.Context, storeID string) ([]StoreBag, error) {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT 
			b.id, b.title, b.description, b.stock, b.original_price, b.final_price, b.status, 
			b.pickup_start, b.pickup_end, b.created_at, b.updated_at,
			i.id, i.name, i.image_key, i.quantity
		FROM bags b
		LEFT JOIN bag_items i ON i.bag_id = b.id
		WHERE b.store_id = $1
		ORDER BY b.created_at DESC, i.name ASC
		`,
		storeID,
	)
	if err != nil {
		return nil, fmt.Errorf("select bags by store id: %v", err)
	}
	defer rows.Close()

	bagsMap := make(map[string]*StoreBag)
	bagsOrder := make([]string, 0)

	for rows.Next() {
		bRow := struct {
			ID            string
			Title         string
			Description   sql.NullString
			Stock         int32
			OriginalPrice int64
			FinalPrice    int64
			Status        string
			PickupStart   time.Time
			PickupEnd     time.Time
			CreatedAt     time.Time
			UpdatedAt     time.Time
		}{}

		iRow := struct {
			ID       sql.NullString
			Name     sql.NullString
			ImageKey sql.NullString
			Quantity sql.NullInt32
		}{}

		if err := rows.Scan(
			&bRow.ID,
			&bRow.Title,
			&bRow.Description,
			&bRow.Stock,
			&bRow.OriginalPrice,
			&bRow.FinalPrice,
			&bRow.Status,
			&bRow.PickupStart,
			&bRow.PickupEnd,
			&bRow.CreatedAt,
			&bRow.UpdatedAt,

			&iRow.ID,
			&iRow.Name,
			&iRow.ImageKey,
			&iRow.Quantity,
		); err != nil {
			return nil, fmt.Errorf("scan bag: %v", err)
		}

		if _, ok := bagsMap[bRow.ID]; !ok {
			pickupWindow := StoreBagPickupWindow{
				Start: bRow.PickupStart,
				End:   bRow.PickupEnd,
			}

			b := &StoreBag{
				ID:            bRow.ID,
				Title:         bRow.Title,
				Stock:         bRow.Stock,
				OriginalPrice: bRow.OriginalPrice,
				FinalPrice:    bRow.FinalPrice,
				Status:        bRow.Status,
				PickupWindow:  pickupWindow,
				CreatedAt:     bRow.CreatedAt,
				UpdatedAt:     bRow.UpdatedAt,
			}

			if bRow.Description.Valid {
				b.Description = &bRow.Description.String
			}

			bagsMap[bRow.ID] = b
			bagsOrder = append(bagsOrder, bRow.ID)
		}

		if iRow.ID.Valid {
			i := StoreBagItem{
				ID:       iRow.ID.String,
				Name:     iRow.Name.String,
				Quantity: iRow.Quantity.Int32,
			}

			if iRow.ImageKey.Valid {
				imgURL := qry.fileStorage.URL(iRow.ImageKey.String)
				i.ImageURL = &imgURL
			}

			bagsMap[bRow.ID].Items = append(bagsMap[bRow.ID].Items, i)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %v", err)
	}

	bags := make([]StoreBag, 0, len(bagsMap))
	for _, id := range bagsOrder {
		bags = append(bags, *bagsMap[id])
	}

	return bags, nil
}
