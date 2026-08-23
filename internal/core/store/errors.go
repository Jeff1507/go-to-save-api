package store

import "errors"

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrDuplicateCNPJ      = errors.New("duplicate cnpj")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("store not found")
	ErrInvalidImageFormat = errors.New("invalid image format")
	ErrIncorrectPassword  = errors.New("incorrect current password")
)
