package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

var ErrContentTypeNotAllowed = errors.New("content type not allowed")

type LocalStorage struct {
	baseURL  string
	basePath string
}

func NewLocalStorage(baseURL, basePath string) (LocalStorage, error) {
	paths := []string{
		ProductsPath,
	}

	for _, path := range paths {
		dir := filepath.Join(basePath, path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return LocalStorage{}, fmt.Errorf("create directory %s: %v", path, err)
		}
	}

	return LocalStorage{
		baseURL:  baseURL,
		basePath: basePath,
	}, nil
}

// TODO: Validar o tamanho do arquivo
func (ls LocalStorage) validateImage(contentType string) error {
	allowed := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/webp": true,
	}

	if !allowed[contentType] {
		return ErrContentTypeNotAllowed
	}

	return nil
}

func (ls LocalStorage) getExt(contentType string) string {
	exts := map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/webp": ".webp",
	}

	if ext, ok := exts[contentType]; ok {
		return ext
	}

	return ".bin"
}

func (ls LocalStorage) Upload(content io.Reader, contentType, path string) (string, error) {
	if err := ls.validateImage(contentType); err != nil {
		if errors.Is(err, ErrContentTypeNotAllowed) {
			return "", fmt.Errorf("validate image: %w", err)
		}

		return "", fmt.Errorf("validate image: %v", err)
	}

	filename := uuid.Generate() + ls.getExt(contentType)
	key := filepath.Join(path, filename)
	fullpath := filepath.Join(ls.basePath, key)

	dst, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("create file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		os.Remove(fullpath)
		return "", fmt.Errorf("copy file content: %v", err)
	}

	return key, nil
}

func (ls LocalStorage) Delete(key string) error {
	if key == "" {
		return nil
	}

	fullpath := filepath.Join(ls.basePath, key)
	if err := os.Remove(fullpath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete image: %v", err)
	}

	return nil
}

func (ls LocalStorage) URL(key string) string {
	return ls.baseURL + "/uploads/" + key
}
