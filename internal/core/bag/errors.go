package bag

import "errors"

var (
	ErrInvalidPickupWindowInterval = errors.New("invalid pickup window interval")
	ErrInvalidOperation            = errors.New("invalid operation")
	ErrNotFound                    = errors.New("bag not found")
	ErrClosed                      = errors.New("closed bag")
	ErrExpired                     = errors.New("expired bag")
)
