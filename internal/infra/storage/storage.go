package storage

import (
	"errors"
	"io"
)

const (
	ProductsPath string = "products"
	AvatarsPath  string = "avatars"
)

var ErrContentTypeNotAllowed = errors.New("content type not allowed")

type Storage interface {
	Upload(content io.Reader, contentType, storagePath string) (string, error)
	Delete(key string) error
	URL(key string) string
}
