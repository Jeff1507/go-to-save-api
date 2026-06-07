package product

import "context"

type Repository interface {
	Save(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, ids []string) ([]*Product, error)
}
