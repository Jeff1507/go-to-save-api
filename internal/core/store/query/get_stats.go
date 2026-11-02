package query

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsOutput struct {
	TodaySales           int64
	ActiveBags           int64
	CompletedOrdersMonth int64
}

type GetStats struct {
	pg    *pgxpool.Pool
	clock clock.Clock
}

func NewGetStats(pg *pgxpool.Pool, clock clock.Clock) *GetStats {
	return &GetStats{
		pg:    pg,
		clock: clock,
	}
}

func (qry *GetStats) Execute(ctx context.Context, storeID string) (StatsOutput, error) {
	now := qry.clock.Now()

	row := qry.pg.QueryRow(
		ctx,
		`
		WITH date_ranges AS (
			SELECT
				date_trunc('day', $2::timestamptz AT TIME ZONE 'America/Sao_Paulo')
					AT TIME ZONE 'America/Sao_Paulo' AS today_start,

				(date_trunc('day', $2::timestamptz AT TIME ZONE 'America/Sao_Paulo')
					+ INTERVAL '1 day') AT TIME ZONE 'America/Sao_Paulo' AS tomorrow_start,

				date_trunc('month', $2::timestamptz AT TIME ZONE 'America/Sao_Paulo')
					AT TIME ZONE 'America/Sao_Paulo' AS month_start,

				(date_trunc('month', $2::timestamptz AT TIME ZONE 'America/Sao_Paulo')
					+ INTERVAL '1 month') AT TIME ZONE 'America/Sao_Paulo' AS next_month_start
		)
		SELECT
			COALESCE((
				SELECT SUM(o.total_price)
				FROM orders o
				CROSS JOIN date_ranges d
				WHERE o.store_id = $1
				  AND o.status = 'COMPLETED'
				  AND o.completed_at >= d.today_start
				  AND o.completed_at < d.tomorrow_start
			), 0)::bigint AS today_sales,

			(
				SELECT COUNT(*)
				FROM bags b
				WHERE b.store_id = $1
				  AND b.status = 'ACTIVE'
				  AND b.stock > 0
				  AND b.pickup_end > $2::timestamptz
			) AS active_bags,

			(
				SELECT COUNT(*)
				FROM orders o
				CROSS JOIN date_ranges d
				WHERE o.store_id = $1
				  AND o.status = 'COMPLETED'
				  AND o.completed_at >= d.month_start
				  AND o.completed_at < d.next_month_start
			) AS completed_orders_month
		`,
		storeID, now,
	)

	var stats StatsOutput
	if err := row.Scan(
		&stats.TodaySales,
		&stats.ActiveBags,
		&stats.CompletedOrdersMonth,
	); err != nil {
		return StatsOutput{}, fmt.Errorf("scan stats: %w", err)
	}

	return stats, nil
}
