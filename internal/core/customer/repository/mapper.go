package repository

import (
	"database/sql"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
)

type customerRow struct {
	ID           string
	Email        string
	FirstName    string
	LastName     string
	PictureKey   sql.NullString
	PasswordHash string
	VerifiedAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func toRow(c *customer.Customer) customerRow {
	var pictureKey sql.NullString
	if c.PictureKey() != "" {
		pictureKey = sql.NullString{String: c.PictureKey(), Valid: true}
	}

	var verifiedAt sql.NullTime
	if c.IsVerified() {
		verifiedAt = sql.NullTime{Time: c.VerifiedAt(), Valid: true}
	}

	return customerRow{
		ID:           c.ID(),
		Email:        c.Email(),
		FirstName:    c.FirstName(),
		LastName:     c.LastName(),
		PictureKey:   pictureKey,
		PasswordHash: c.PasswordHash(),
		VerifiedAt:   verifiedAt,
		CreatedAt:    c.CreatedAt(),
		UpdatedAt:    c.UpdatedAt(),
	}
}

func toEntity(r customerRow) *customer.Customer {
	var pictureKey string
	if r.PictureKey.Valid {
		pictureKey = r.PictureKey.String
	}

	var verifiedAt time.Time
	if r.VerifiedAt.Valid {
		verifiedAt = r.VerifiedAt.Time
	}

	return customer.From(
		r.ID,
		r.Email,
		r.FirstName,
		r.LastName,
		pictureKey,
		r.PasswordHash,
		verifiedAt,
		r.CreatedAt,
		r.UpdatedAt,
	)
}
