package store

import "errors"

var (
	ErrEmailAlreadyInUse  = errors.New("email is already in use")
	ErrDuplicateCNPJ      = errors.New("duplicate cnpj")
	ErrNotFound           = errors.New("store not found")
	ErrInvalidImageFormat = errors.New("invalid image format")
)
