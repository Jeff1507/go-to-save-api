package repository

import (
	"context"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
)

type InMemRepo struct {
	presets map[string]preset.Preset
}

func NewInMemRepo() *InMemRepo {
	return &InMemRepo{
		presets: make(map[string]preset.Preset),
	}
}

func (repo *InMemRepo) isPresetExists(p *preset.Preset) bool {
	for _, pset := range repo.presets {
		if pset.ID() == p.ID() {
			continue
		}

		if pset.StoreID() == p.StoreID() && pset.Name() == p.Name() {
			return true
		}
	}

	return false
}

func (repo *InMemRepo) Save(ctx context.Context, p *preset.Preset) error {
	if _, ok := repo.presets[p.ID()]; ok {
		return preset.ErrAlreadyExists
	}

	if repo.isPresetExists(p) {
		return preset.ErrAlreadyExists
	}

	repo.presets[p.ID()] = *p

	return nil
}

func (repo *InMemRepo) Update(ctx context.Context, p *preset.Preset) error {
	if repo.isPresetExists(p) {
		return preset.ErrAlreadyExists
	}

	repo.presets[p.ID()] = *p

	return nil
}

func (repo *InMemRepo) Delete(ctx context.Context, id string) error {
	delete(repo.presets, id)

	return nil
}

func (repo *InMemRepo) Find(ctx context.Context, id string) (*preset.Preset, error) {
	if pset, ok := repo.presets[id]; ok {
		return &pset, nil
	}

	return nil, preset.ErrNotFound
}
