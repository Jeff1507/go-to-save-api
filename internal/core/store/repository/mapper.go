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
	AvatarKey    sql.NullString
	Description  sql.NullString
	PasswordHash string
	RatingSum    int32
	RatingCount  int32
	Address      addressRow
	VerifiedAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func toRow(s *store.Store) storeRow {
	var avatarKey sql.NullString
	if s.AvatarKey() != "" {
		avatarKey = sql.NullString{String: s.AvatarKey(), Valid: true}
	}

	var description sql.NullString
	if s.Description() != "" {
		description = sql.NullString{String: s.Description(), Valid: true}
	}

	var verifiedAt sql.NullTime
	if s.IsVerified() {
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
		AvatarKey:    avatarKey,
		Description:  description,
		PasswordHash: s.PasswordHash(),
		RatingSum:    s.Rating().Sum(),
		RatingCount:  s.Rating().Count(),
		Address:      address,
		VerifiedAt:   verifiedAt,
		CreatedAt:    s.CreatedAt(),
		UpdatedAt:    s.UpdatedAt(),
	}
}

func toEntity(r storeRow) *store.Store {
	rating := store.NewRating(r.RatingSum, r.RatingCount)

	address, _ := store.NewAddress(
		r.Address.PostalCode,
		r.Address.Street,
		r.Address.Number,
		r.Address.Neighborhood,
		r.Address.Complement.String,
		r.Address.City,
		r.Address.State,
		r.Address.Latitude,
		r.Address.Longitude,
	)

	return store.From(
		r.ID,
		r.Email,
		r.Name,
		r.CNPJ,
		r.AvatarKey.String,
		r.Description.String,
		r.PasswordHash,
		rating,
		address,
		r.VerifiedAt.Time,
		r.CreatedAt,
		r.UpdatedAt,
	)
}
