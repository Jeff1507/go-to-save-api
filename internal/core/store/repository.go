package store

import "context"

type Repository interface {
	Save(ctx context.Context, s *Store) error
	Update(ctx context.Context, s *Store) error
	Exists(ctx context.Context, email string) (bool, error)
	Find(ctx context.Context, id string) (*Store, error)
	FindByEmail(ctx context.Context, email string) (*Store, error)
}
