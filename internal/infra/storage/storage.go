package storage

import "io"

const (
	ProductsPath string = "products"
)

type Storage interface {
	Upload(content io.Reader, contentType, path string) (string, error)
	Delete(key string) error
	URL(key string) string
}
