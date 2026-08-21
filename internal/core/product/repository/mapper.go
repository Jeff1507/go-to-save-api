package repository

import (
	"database/sql"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
)

type productRow struct {
	ID        string
	StoreID   string
	Name      string
	ImageKey  sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toRow(p *product.Product) productRow {
	var imageKey sql.NullString
	if p.ImageKey() != "" {
		imageKey = sql.NullString{String: p.ImageKey(), Valid: true}
	}

	return productRow{
		ID:        p.ID(),
		StoreID:   p.StoreID(),
		Name:      p.Name(),
		ImageKey:  imageKey,
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}

func toEntity(r productRow) *product.Product {
	return product.From(r.ID, r.StoreID, r.Name, r.ImageKey.String, r.CreatedAt, r.UpdatedAt)
}
