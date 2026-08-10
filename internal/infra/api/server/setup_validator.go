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

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}

		return strings.Split(name, ",")[0]
	})

	v.RegisterValidation("notblank", validation.NotBlank)
	v.RegisterValidation("time", validation.Time)
}
