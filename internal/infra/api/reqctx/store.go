package reqctx

import "context"

type storeKey struct{}

func WithStore(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, storeKey{}, id)
}

func StoreID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(storeKey{}).(string)
	return id, ok
}
