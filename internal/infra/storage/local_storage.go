package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (LocalStorage, error) {
	paths := []string{AvatarsPath, ProductsPath}

	for _, path := range paths {
		dir := filepath.Join(basePath, path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return LocalStorage{}, fmt.Errorf("create directory %s: %w", path, err)
		}
	}

	return LocalStorage{
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

func (ls LocalStorage) Upload(content io.Reader, contentType, storagePath string) (string, error) {
	if err := ls.validateImage(contentType); err != nil {
		return "", fmt.Errorf("validate image: %w", err)
	}

	filename := uuid.Generate() + ls.getExt(contentType)
	key := path.Join(storagePath, filename)
	fullpath := filepath.Join(ls.basePath, storagePath, filename)

	dst, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		os.Remove(fullpath)
		return "", fmt.Errorf("copy file content: %w", err)
	}

	return key, nil
}

func (ls LocalStorage) Delete(key string) error {
	if key == "" {
		return nil
	}

	fullpath := filepath.Join(ls.basePath, key)
	if err := os.Remove(fullpath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete image: %w", err)
	}

	return nil
}

func (ls LocalStorage) URL(key string) string {
	return "/uploads/" + key
}
