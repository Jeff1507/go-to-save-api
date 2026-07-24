package store

import (
	"errors"
	"time"
)

type Store struct {
	id           string
	email        string
	name         string
	cnpj         string
	pictureKey   string // Opcional
	description  string // Opcional
	passwordHash string
	address      Address
	verifiedAt   time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func New(
	id string,
	email string,
	name string,
	cnpj string,
	passwordHash string,
	address Address,
	now time.Time,
) (*Store, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if cnpj == "" {
		return nil, errors.New("cnpj is required")
	}

	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}

	if now.IsZero() {
		return nil, errors.New("now is not set")
	}

	return &Store{
		id:           id,
		email:        email,
		name:         name,
		cnpj:         cnpj,
		passwordHash: passwordHash,
		address:      address,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func From(
	id string,
	email string,
	name string,
	cnpj string,
	pictureKey string,
	description string,
	passwordHash string,
	address Address,
	verifiedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Store {
	return &Store{
		id:           id,
		email:        email,
		name:         name,
		cnpj:         cnpj,
		pictureKey:   pictureKey,
		description:  description,
		passwordHash: passwordHash,
		address:      address,
		verifiedAt:   verifiedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (s Store) IsVerified() bool {
	return !s.verifiedAt.IsZero()
}

func (s *Store) Verify(now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.verifiedAt = now
	s.updatedAt = now

	return nil
}

func (s *Store) ChangeEmail(email string, now time.Time) error {
	if email == "" {
		return errors.New("email is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.email = email
	s.updatedAt = now

	return nil
}

func (s *Store) ChangeName(name string, now time.Time) error {
	if name == "" {
		return errors.New("name is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.name = name
	s.updatedAt = now

	return nil
}

func (s *Store) ChangeCNPJ(cnpj string, now time.Time) error {
	if cnpj == "" {
		return errors.New("cnpj is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.cnpj = cnpj
	s.updatedAt = now

	return nil
}

func (s *Store) ChangePicture(pictureKey string, now time.Time) error {
	if pictureKey == "" {
		return nil
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.pictureKey = pictureKey
	s.updatedAt = now

	return nil
}

func (s *Store) ChangeDescription(description string, now time.Time) error {
	if description == "" {
		return nil
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.description = description
	s.updatedAt = now

	return nil
}

func (s *Store) ChangePasswordHash(passwordHash string, now time.Time) error {
	if passwordHash == "" {
		return errors.New("password hash is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.passwordHash = passwordHash
	s.updatedAt = now

	return nil
}

func (s *Store) ChangeAddress(address Address, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.address = address
	s.updatedAt = now

	return nil
}

func (s Store) ID() string {
	return s.id
}

func (s Store) Email() string {
	return s.email
}

func (s Store) Name() string {
	return s.name
}

func (s Store) CNPJ() string {
	return s.cnpj
}

func (s Store) PictureKey() string {
	return s.pictureKey
}

func (s Store) Description() string {
	return s.description
}

func (s Store) PasswordHash() string {
	return s.passwordHash
}

func (s Store) Address() Address {
	return s.address
}

func (s Store) VerifiedAt() time.Time {
	return s.verifiedAt
}

func (s Store) CreatedAt() time.Time {
	return s.createdAt
}

func (s Store) UpdatedAt() time.Time {
	return s.updatedAt
}
