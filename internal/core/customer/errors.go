package customer

import "errors"

var (
	ErrEmailAlreadyInUse  = errors.New("email is already in use")
	ErrNotFound           = errors.New("customer not found")
	ErrInvalidImageFormat = errors.New("invalid image format")
)
