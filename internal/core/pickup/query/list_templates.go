package query

import (
	"context"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateOutput struct {
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

func (qry *ListTemplates) Execute(ctx context.Context, storeID string) ([]TemplateOutput, error) {
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

	templates := make([]TemplateOutput, 0)
	for rows.Next() {
		var (
			templ TemplateOutput
			start time.Time
			end   time.Time
		)
		if err := rows.Scan(
			&templ.ID,
			&templ.Name,
			&start,
			&end,
		); err != nil {
			return nil, fmt.Errorf("scan pickup template: %w", err)
		}

		templ.Start = pickup.FormatTime(start.Hour(), start.Minute())
		templ.End = pickup.FormatTime(end.Hour(), end.Minute())

		templates = append(templates, templ)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return templates, nil
}
