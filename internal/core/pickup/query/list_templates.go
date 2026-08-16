package query

import (
	"context"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Template struct {
	ID    string
	Name  string
	Start string
	End   string
}

type ListTemplates struct {
	pg *pgxpool.Pool
}

func NewListTemplates(pg *pgxpool.Pool) *ListTemplates {
	return &ListTemplates{
		pg: pg,
	}
}

func (qry *ListTemplates) Execute(ctx context.Context, storeID string) ([]Template, error) {
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
		return nil, fmt.Errorf("select pickup templates: %w", err)
	}
	defer rows.Close()

	templates := make([]Template, 0)
	for rows.Next() {
		var (
			template Template
			start    time.Time
			end      time.Time
		)
		if err := rows.Scan(
			&template.ID,
			&template.Name,
			&start,
			&end,
		); err != nil {
			return nil, fmt.Errorf("scan pickup template: %w", err)
		}

		template.Start = pickup.FormatTime(start.Hour(), start.Minute())
		template.End = pickup.FormatTime(end.Hour(), end.Minute())

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return templates, nil
}
