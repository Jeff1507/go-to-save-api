package query

import (
	"context"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PresetOutput struct {
	ID    string
	Name  string
	Start string
	End   string
}

type ListPresets struct {
	pg *pgxpool.Pool
}

func NewListPresets(pg *pgxpool.Pool) *ListPresets {
	return &ListPresets{
		pg: pg,
	}
}

func (qry *ListPresets) Execute(ctx context.Context, storeID string) ([]PresetOutput, error) {
	rows, err := qry.pg.Query(
		ctx,
		`
		SELECT id, name, start_time, end_time
		FROM pickup_presets
		WHERE store_id = $1
		ORDER BY name ASC
		`,
		storeID,
	)
	if err != nil {
		return nil, fmt.Errorf("select presets: %w", err)
	}
	defer rows.Close()

	presets := make([]PresetOutput, 0)
	for rows.Next() {
		var (
			pset  PresetOutput
			start time.Time
			end   time.Time
		)
		if err := rows.Scan(
			&pset.ID,
			&pset.Name,
			&start,
			&end,
		); err != nil {
			return nil, fmt.Errorf("scan presets: %w", err)
		}

		pset.Start = pickup.FormatTime(start.Hour(), start.Minute())
		pset.End = pickup.FormatTime(end.Hour(), end.Minute())

		presets = append(presets, pset)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return presets, nil
}
