package reqctx

import "context"

type customerKey struct{}

func WithCustomer(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, customerKey{}, id)
}

func CustomerID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(customerKey{}).(string)
	return id, ok
}
