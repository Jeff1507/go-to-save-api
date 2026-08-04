package request

import "context"

type customerKey struct{}

func WithCustomerID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, customerKey{}, id)
}

func CustomerID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(customerKey{}).(string)
	return id, ok
}

type storeKey struct{}

func WithStoreID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, storeKey{}, id)
}

func StoreID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(storeKey{}).(string)
	return id, ok
}
