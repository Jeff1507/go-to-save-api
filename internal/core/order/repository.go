package order

import "context"

type Repository interface {
	Save(ctx context.Context, o *Order) error
	Update(ctx context.Context, o *Order) error
	Has(ctx context.Context, bagID string) (bool, error)
	Find(ctx context.Context, id string) (*Order, error)
}
