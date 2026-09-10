package bag

import "context"

type Repository interface {
	Save(ctx context.Context, b *Bag) error
	Update(ctx context.Context, b *Bag) error
	Find(ctx context.Context, id string) (*Bag, error)
	FindAndUpdate(ctx context.Context, id string, fn func(b *Bag) error) error
}
