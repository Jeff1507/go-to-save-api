package product

import (
	"errors"
	"time"
)

type Product struct {
	id        string
	storeID   string
	name      string
	imageKey  string // Opcional
	createdAt time.Time
	updatedAt time.Time
}

func New(id, storeID, name, imageKey string, now time.Time) (*Product, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if storeID == "" {
		return nil, errors.New("store id is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if now.IsZero() {
		return nil, errors.New("now is not set")
	}

	return &Product{
		id:        id,
		storeID:   storeID,
		name:      name,
		imageKey:  imageKey,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func From(id, storeID, name, imageKey string, createdAt, updatedAt time.Time) *Product {
	return &Product{
		id:        id,
		storeID:   storeID,
		name:      name,
		imageKey:  imageKey,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (p *Product) ChangeName(name string, now time.Time) error {
	if name == "" {
		return errors.New("name is required")
	}

	if now.IsZero() {
		return errors.New("now is not set")
	}

	p.name = name
	p.updatedAt = now

	return nil
}

func (p *Product) ChangeImage(imageKey string, now time.Time) error {
	if now.IsZero() {
		return errors.New("now is not set")
	}

	p.imageKey = imageKey
	p.updatedAt = now

	return nil
}

func (p *Product) RemoveImage(now time.Time) error {
	return p.ChangeImage("", now)
}

func (p Product) ID() string {
	return p.id
}

func (p Product) StoreID() string {
	return p.storeID
}

func (p Product) Name() string {
	return p.name
}

func (p Product) ImageKey() string {
	return p.imageKey
}

func (p Product) CreatedAt() time.Time {
	return p.createdAt
}

func (p Product) UpdatedAt() time.Time {
	return p.updatedAt
}
