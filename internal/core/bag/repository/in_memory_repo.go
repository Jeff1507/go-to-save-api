package repository

import (
	"context"
	"errors"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
)

type InMemRepo struct {
	bags map[string]bag.Bag
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		bags: make(map[string]bag.Bag),
	}
}

func (repo *InMemRepo) Save(ctx context.Context, b *bag.Bag) error {
	if _, ok := repo.bags[b.ID()]; ok {
		return errors.New("bag already created")
	}

	repo.bags[b.ID()] = *b

	return nil
}

func (repo *InMemRepo) Update(ctx context.Context, b *bag.Bag) error {
	if _, ok := repo.bags[b.ID()]; !ok {
		return bag.ErrNotFound
	}

	repo.bags[b.ID()] = *b

	return nil
}

func (repo *InMemRepo) Find(ctx context.Context, id string) (*bag.Bag, error) {
	if b, ok := repo.bags[id]; ok {
		return &b, nil
	}

	return nil, bag.ErrNotFound
}
