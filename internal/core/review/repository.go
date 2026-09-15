package review

import "context"

type Repository interface {
	Save(ctx context.Context, r *Review) error
}
