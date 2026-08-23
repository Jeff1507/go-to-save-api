package storage

import (
	"errors"
	"io"
)

const (
	AvatarsPath  string = "avatars"
	ProductsPath string = "products"
)

var ErrContentTypeNotAllowed = errors.New("content type not allowed")

type Storage interface {
	Upload(content io.Reader, contentType, storagePath string) (string, error)
	Delete(key string) error
	URL(key string) string
}
