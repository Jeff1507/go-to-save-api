package server

import (
	"reflect"
	"strings"

	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func setupValidator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("failed to get validator engine")
	}

	v.RegisterTagNameFunc(getTagName)

	v.RegisterValidation("notblank", validation.NotBlank)
	v.RegisterValidation("time", validation.Time)
}

func getTagName(field reflect.StructField) string {
	tags := []string{"json", "form", "uri"}

	for _, tag := range tags {
		if name, _, _ := strings.Cut(field.Tag.Get(tag), ","); name != "" {
			return name
		}
	}

	return field.Name
}
