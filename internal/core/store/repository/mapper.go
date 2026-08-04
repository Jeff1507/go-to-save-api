package repository

import (
	"database/sql"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
)

type addressRow struct {
	PostalCode   string
	Street       string
	Number       string
	Neighborhood string
	Complement   sql.NullString
	City         string
	State        string
	Latitude     float64
	Longitude    float64
}

type storeRow struct {
	ID           string
	Email        string
	Name         string
	CNPJ         string
	PictureKey   sql.NullString
	Description  sql.NullString
	PasswordHash string
	Address      addressRow
	VerifiedAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func toRow(s *store.Store) storeRow {
	var pictureKey sql.NullString
	if s.PictureKey() != "" {
		pictureKey = sql.NullString{String: s.PictureKey(), Valid: true}
	}

	var description sql.NullString
	if s.Description() != "" {
		description = sql.NullString{String: s.Description(), Valid: true}
	}

	var verifiedAt sql.NullTime
	if !s.VerifiedAt().IsZero() {
		verifiedAt = sql.NullTime{Time: s.VerifiedAt(), Valid: true}
	}

	var addressComplement sql.NullString
	if s.Address().Complement() != "" {
		addressComplement = sql.NullString{String: s.Address().Complement(), Valid: true}
	}

	address := addressRow{
		PostalCode:   s.Address().PostalCode(),
		Street:       s.Address().Street(),
		Number:       s.Address().Number(),
		Neighborhood: s.Address().Neighborhood(),
		Complement:   addressComplement,
		City:         s.Address().City(),
		State:        s.Address().State(),
		Latitude:     s.Address().Latitude(),
		Longitude:    s.Address().Longitude(),
	}

	return storeRow{
		ID:           s.ID(),
		Email:        s.Email(),
		Name:         s.Name(),
		CNPJ:         s.CNPJ(),
		PictureKey:   pictureKey,
		Description:  description,
		PasswordHash: s.PasswordHash(),
		Address:      address,
		VerifiedAt:   verifiedAt,
		CreatedAt:    s.CreatedAt(),
		UpdatedAt:    s.UpdatedAt(),
	}
}
