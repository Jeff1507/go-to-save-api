package preset

import "context"

type Repository interface {
	Save(ctx context.Context, p *Preset) error
	Update(ctx context.Context, p *Preset) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (*Preset, error)
}
