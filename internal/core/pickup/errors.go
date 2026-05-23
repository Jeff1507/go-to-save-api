package pickup

import "errors"

var (
	ErrInvalidInterval       = errors.New("invalid pickup interval")
	ErrTemplateAlreadyExists = errors.New("pickup template already exists")
	ErrTemplateNotFound      = errors.New("pickup template not found")
)
