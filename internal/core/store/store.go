package store

import (
	"errors"
	"fmt"
	"time"
)

type Store struct {
	id           string
	email        string
	name         string
	cnpj         string
	avatarKey    string // Opcional
	description  string // Opcional
	passwordHash string
	address      Address
	rating       Rating
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

	if address == (Address{}) {
		return nil, errors.New("address is required")
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
	avatarKey string,
	description string,
	passwordHash string,
	rating Rating,
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
		avatarKey:    avatarKey,
		description:  description,
		passwordHash: passwordHash,
		rating:       rating,
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

func (s *Store) AddRating(rating int16) error {
	if err := s.rating.Add(rating); err != nil {
		return fmt.Errorf("add rating: %w", err)
	}

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
	s.verifiedAt = time.Time{}
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

func (s *Store) ChangeAvatar(avatarKey string, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.avatarKey = avatarKey
	s.updatedAt = now

	return nil
}

func (s *Store) RemoveAvatar(now time.Time) error {
	return s.ChangeAvatar("", now)
}

func (s *Store) ChangeDescription(description string, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	s.description = description
	s.updatedAt = now

	return nil
}

func (s *Store) ChangePassword(passwordHash string, now time.Time) error {
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
	if address == (Address{}) {
		return errors.New("address is required")
	}

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

func (s Store) AvatarKey() string {
	return s.avatarKey
}

func (s Store) Description() string {
	return s.description
}

func (s Store) PasswordHash() string {
	return s.passwordHash
}

func (s Store) Rating() Rating {
	return s.rating
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
