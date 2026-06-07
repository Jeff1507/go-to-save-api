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
	var imgKey sql.NullString
	if p.ImageKey() != "" {
		imgKey = sql.NullString{String: p.ImageKey(), Valid: true}
	}

	return productRow{
		ID:        p.ID(),
		StoreID:   p.StoreID(),
		Name:      p.Name(),
		ImageKey:  imgKey,
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}

func toEntity(r productRow) *product.Product {
	var imgKey string
	if r.ImageKey.Valid {
		imgKey = r.ImageKey.String
	}

	return product.From(r.ID, r.StoreID, r.Name, imgKey, r.CreatedAt, r.UpdatedAt)
}
