package bag

import "errors"

var (
	ErrNotFound          = errors.New("bag not found")
	ErrClosed            = errors.New("closed bag")
	ErrExpired           = errors.New("expired bag")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrDuplicateItems    = errors.New("duplicate items")
	ErrProductNotFound   = errors.New("product not found")
)
