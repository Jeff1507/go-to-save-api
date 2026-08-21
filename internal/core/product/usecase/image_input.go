package usecase

import "io"

type ImageInput struct {
	Content     io.Reader
	ContentType string
}
