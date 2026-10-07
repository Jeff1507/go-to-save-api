package order

import "context"

type Repository interface {
	Save(ctx context.Context, o *Order) error
	Update(ctx context.Context, o *Order) error
	HasPendingByBag(ctx context.Context, bagID string) (bool, error)
	HasPendingByCustomer(ctx context.Context, customerID string) (bool, error)
	Find(ctx context.Context, id string) (*Order, error)
}
