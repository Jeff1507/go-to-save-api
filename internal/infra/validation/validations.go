package validation

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var NotBlank validator.Func = func(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

var Time validator.Func = func(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04", fl.Field().String())
	return err == nil
}
