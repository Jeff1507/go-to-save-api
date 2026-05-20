package product

import "errors"

var (
	ErrAlreadyExists = errors.New("product already exists")
	ErrNotFound      = errors.New("product not found")
)
