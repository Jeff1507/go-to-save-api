package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}

		return name
	})

	v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	v.RegisterValidation("time", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("15:04", fl.Field().String())
		return err == nil
	})

	return v
}

func UUID(v *validator.Validate, errs Errors, val, field string) {
	if err := v.Var(val, "required,uuid"); err != nil {
		if vErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range vErr {
				errs.Add(field, translateTag(e))
			}
		}
	}
}

func Struct(v *validator.Validate, errs Errors, s any) {
	if err := v.Struct(s); err != nil {
		if vErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range vErrs {
				field := formatNamespace(e.Namespace())
				errs.Add(field, translateTag(e))
			}
		}
	}
}

func formatNamespace(ns string) string {
	return ns[strings.Index(ns, ".")+1:]
}

func translateTag(fErr validator.FieldError) string {
	switch fErr.Tag() {
	case "required":
		return "obrigatório"

	case "notblank":
		return "deve ser preenchido"

	case "uuid":
		return "deve ser um uuid"

	case "email":
		return "deve ser um e-mail válido"

	case "time":
		return "deve estar no formato HH:MM"

	case "min":
		if fErr.Kind() == reflect.Slice {
			return fmt.Sprintf("deve ter no mínimo %s elemento(s)", fErr.Param())
		}
		return fmt.Sprintf("deve ter no mínimo %s caracteres", fErr.Param())

	case "max":
		if fErr.Kind() == reflect.Slice {
			return fmt.Sprintf("deve ter no máximo %s elemento(s)", fErr.Param())
		}
		return fmt.Sprintf("deve ter no máximo %s caracteres", fErr.Param())

	case "gt":
		return fmt.Sprintf("deve ser maior que %s", fErr.Param())

	case "gte":
		return fmt.Sprintf("deve ser maior ou igual a %s", fErr.Param())

	case "lt":
		return fmt.Sprintf("deve ser menor que %s", fErr.Param())

	case "lte":
		return fmt.Sprintf("deve ser menor ou igual a %s", fErr.Param())

	case "gtfield":
		return fmt.Sprintf("deve ser posterior ao %s", strings.ToLower(fErr.Param()))

	default:
		return fErr.Error()
	}
}
