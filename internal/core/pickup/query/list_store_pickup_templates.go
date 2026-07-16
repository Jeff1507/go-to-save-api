package query

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StorePickupTemplate struct {
	ID    string
	Name  string
	Start string
	End   string
}

type ListStorePickupTemplates struct {
	pg *pgxpool.Pool
}

func NewListStorePickupTemplates(pg *pgxpool.Pool) *ListStorePickupTemplates {
	return &ListStorePickupTemplates{
		pg: pg,
	}
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func (qry *ListStorePickupTemplates) Execute(
	ctx context.Context,
	storeID string,
) ([]StorePickupTemplate, error) {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT id, name, start_time, end_time
		FROM pickup_templates
		WHERE store_id = $1
		ORDER BY name ASC
		`,
		storeID,
	)
	if err != nil {
		return nil, fmt.Errorf("select pickup templates: %v", err)
	}
	defer rows.Close()

	templates := make([]StorePickupTemplate, 0)
	for rows.Next() {
		var start, end time.Time
		var t StorePickupTemplate
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&start,
			&end,
		); err != nil {
			return nil, fmt.Errorf("scan pickup template: %v", err)
		}

		t.Start = formatTime(start)
		t.End = formatTime(end)

		templates = append(templates, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %v", err)
	}

	return templates, nil
}
