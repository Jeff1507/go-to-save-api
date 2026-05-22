package storage

import (
	"fmt"
	"io"
	"net/http"
)

func GetContentType(file io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)

	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("read buffer: %v", err)
	}

	contentType := http.DetectContentType(buffer)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("reset cursor: %v", err)
	}

	return contentType, nil
}
