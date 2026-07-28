package customer

import "context"

type Repository interface {
	Save(ctx context.Context, c *Customer) error
	Update(ctx context.Context, c *Customer) error
	Find(ctx context.Context, id string) (*Customer, error)
	FindByEmail(ctx context.Context, email string) (*Customer, error)
}
