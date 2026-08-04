package customer

import "errors"

var (
	ErrEmailAlreadyInUse  = errors.New("email is already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("customer not found")
	ErrInvalidImageFormat = errors.New("invalid image format")
)
