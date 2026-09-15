package repository

import (
	"database/sql"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/review"
)

type reviewRow struct {
	ID         string
	OrderID    string
	Rating     int16
	Comment    sql.NullString
	ReviewedAt time.Time
}

func toRow(r *review.Review) reviewRow {
	var comment sql.NullString
	if r.Comment() != "" {
		comment = sql.NullString{String: r.Comment(), Valid: true}
	}

	return reviewRow{
		ID:         r.ID(),
		OrderID:    r.OrderID(),
		Rating:     r.Rating(),
		Comment:    comment,
		ReviewedAt: r.ReviewedAt(),
	}
}
