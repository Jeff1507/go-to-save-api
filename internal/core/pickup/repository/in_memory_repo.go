package repository

import (
	"context"

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

func (repo *InMemRepo) isTemplateExists(t *pickup.Template) bool {
	for _, templ := range repo.templates {
		if templ.ID() == t.ID() {
			continue
		}

		if templ.StoreID() == t.StoreID() && templ.Name() == t.Name() {
			return true
		}
	}

	return false
}

func (repo *InMemRepo) Save(ctx context.Context, t *pickup.Template) error {
	if _, ok := repo.templates[t.ID()]; ok {
		return pickup.ErrTemplateAlreadyExists
	}

	if repo.isTemplateExists(t) {
		return pickup.ErrTemplateAlreadyExists
	}

	repo.templates[t.ID()] = *t

	return nil
}

func (repo *InMemRepo) Update(ctx context.Context, t *pickup.Template) error {
	if repo.isTemplateExists(t) {
		return pickup.ErrTemplateAlreadyExists
	}

	repo.templates[t.ID()] = *t

	return nil
}

func (repo *InMemRepo) Delete(ctx context.Context, id string) error {
	delete(repo.templates, id)

	return nil
}

func (repo *InMemRepo) Find(ctx context.Context, id string) (*pickup.Template, error) {
	if templ, ok := repo.templates[id]; ok {
		return &templ, nil
	}

	return nil, pickup.ErrTemplateNotFound
}
