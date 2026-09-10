package order

import "errors"

var (
	ErrInvalidState     = errors.New("invalid state")
	ErrDuplicateItems   = errors.New("duplicate items")
	ErrInvalidCode      = errors.New("invalid confirmation code")
	ErrNotFound         = errors.New("order not found")
	ErrCannotPlaceOrder = errors.New("can not place order")
	ErrItemNotFound     = errors.New("item not found")
)
