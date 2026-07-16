package pickup

import "context"

type Repository interface {
	Save(ctx context.Context, t *Template) error
	Update(ctx context.Context, t *Template) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (*Template, error)
}
