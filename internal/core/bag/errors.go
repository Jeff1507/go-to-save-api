package bag

import "errors"

var (
	ErrNotFound          = errors.New("bag not found")
	ErrClosed            = errors.New("closed bag")
	ErrExpired           = errors.New("expired bag")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrDuplicateItems    = errors.New("duplicate items")

	ErrItemNotFound    = errors.New("item not found")
	ErrHasPendingOrder = errors.New("bag has a pending order and cannot be edited")
	ErrZeroDelta       = errors.New("delta must not be zero")
)
