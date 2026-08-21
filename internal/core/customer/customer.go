package customer

import (
	"errors"
	"time"
)

type Customer struct {
	id           string
	email        string
	firstName    string
	lastName     string
	avatarKey    string
	passwordHash string
	verifiedAt   time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func New(id, email, firstName, lastName, passwordHash string, now time.Time) (*Customer, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if firstName == "" {
		return nil, errors.New("first name is required")
	}

	if lastName == "" {
		return nil, errors.New("last name is required")
	}

	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}

	if now.IsZero() {
		return nil, errors.New("now is not set")
	}

	return &Customer{
		id:           id,
		email:        email,
		firstName:    firstName,
		lastName:     lastName,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func From(
	id string,
	email string,
	firstName string,
	lastName string,
	avatarKey string,
	passwordHash string,
	verifiedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Customer {
	return &Customer{
		id:           id,
		email:        email,
		firstName:    firstName,
		lastName:     lastName,
		avatarKey:    avatarKey,
		passwordHash: passwordHash,
		verifiedAt:   verifiedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (c Customer) IsVerified() bool {
	return !c.verifiedAt.IsZero()
}

func (c *Customer) Verify(now time.Time) error {
	if c.IsVerified() {
		return nil
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	c.verifiedAt = now
	c.updatedAt = now

	return nil
}

func (c *Customer) ChangeEmail(email string, now time.Time) error {
	if email == "" {
		return errors.New("email is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	c.email = email
	c.verifiedAt = time.Time{}
	c.updatedAt = now

	return nil
}

func (c *Customer) ChangeName(firstName, lastName string, now time.Time) error {
	if firstName == "" {
		return errors.New("first name is required")
	}

	if lastName == "" {
		return errors.New("last name is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	c.firstName = firstName
	c.lastName = lastName
	c.updatedAt = now

	return nil
}

func (c *Customer) ChangeAvatar(avatarKey string, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	c.avatarKey = avatarKey
	c.updatedAt = now

	return nil
}

func (c *Customer) RemoveAvatar(now time.Time) error {
	return c.ChangeAvatar("", now)
}

func (c *Customer) ChangePassword(passwordHash string, now time.Time) error {
	if passwordHash == "" {
		return errors.New("password hash is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	c.passwordHash = passwordHash
	c.updatedAt = now

	return nil
}

func (c Customer) ID() string {
	return c.id
}

func (c Customer) Email() string {
	return c.email
}

func (c Customer) FirstName() string {
	return c.firstName
}

func (c Customer) LastName() string {
	return c.lastName
}

func (c Customer) FullName() string {
	return c.firstName + " " + c.lastName
}

func (c Customer) AvatarKey() string {
	return c.avatarKey
}

func (c Customer) PasswordHash() string {
	return c.passwordHash
}

func (c Customer) VerifiedAt() time.Time {
	return c.verifiedAt
}

func (c Customer) CreatedAt() time.Time {
	return c.createdAt
}

func (c Customer) UpdatedAt() time.Time {
	return c.updatedAt
}
