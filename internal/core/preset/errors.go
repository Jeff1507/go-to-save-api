package preset

import "errors"

var (
	ErrInvalidInterval = errors.New("invalid preset interval")
	ErrAlreadyExists   = errors.New("preset already exists")
	ErrNotFound        = errors.New("preset not found")
)
