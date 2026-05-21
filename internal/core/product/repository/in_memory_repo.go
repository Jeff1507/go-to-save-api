package repository

import (
	"context"
	"errors"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
)

type InMemRepo struct {
	products map[string]product.Product
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		products: make(map[string]product.Product),
	}
}

func (repo *InMemRepo) Save(ctx context.Context, p *product.Product) error {
	if _, ok := repo.products[p.ID()]; ok {
		return errors.New("product already created")
	}

	repo.products[p.ID()] = *p

	return nil
}

func (repo *InMemRepo) Update(ctx context.Context, p *product.Product) error {
	repo.products[p.ID()] = *p

	return nil
}

func (repo *InMemRepo) Delete(ctx context.Context, id string) error {
	delete(repo.products, id)

	return nil
}

func (repo *InMemRepo) Exists(ctx context.Context, storeID, name string) (bool, error) {
	for _, prod := range repo.products {
		if prod.StoreID() == storeID && prod.Name() == name {
			return true, nil
		}
	}

	return false, nil
}

func (repo *InMemRepo) Find(ctx context.Context, id string) (*product.Product, error) {
	if prod, ok := repo.products[id]; ok {
		return &prod, nil
	}

	return nil, product.ErrNotFound
}

func (repo *InMemRepo) List(ctx context.Context, ids []string) ([]*product.Product, error) {
	products := make([]*product.Product, 0, len(ids))

	for _, id := range ids {
		if prod, ok := repo.products[id]; ok {
			products = append(products, &prod)
		}
	}

	return products, nil
}
