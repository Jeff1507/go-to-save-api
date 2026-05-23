package repository

import (
	"context"
	"errors"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type InMemRepo struct {
	templates map[string]pickup.Template
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		templates: make(map[string]pickup.Template),
	}
}

func (repo *InMemRepo) Save(ctx context.Context, t *pickup.Template) error {
	if _, ok := repo.templates[t.ID()]; ok {
		return errors.New("pickup template already created")
	}

	repo.templates[t.ID()] = *t

	return nil
}

func (repo *InMemRepo) Update(ctx context.Context, t *pickup.Template) error {
	repo.templates[t.ID()] = *t

	return nil
}

func (repo *InMemRepo) Delete(ctx context.Context, id string) error {
	delete(repo.templates, id)

	return nil
}

func (repo *InMemRepo) Exists(ctx context.Context, storeID, name string) (bool, error) {
	for _, templ := range repo.templates {
		if templ.StoreID() == storeID && templ.Name() == name {
			return true, nil
		}
	}

	return false, nil
}

func (repo *InMemRepo) Find(ctx context.Context, id string) (*pickup.Template, error) {
	if templ, ok := repo.templates[id]; ok {
		return &templ, nil
	}

	return nil, pickup.ErrTemplateNotFound
}
