package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(vErrs validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, e := range vErrs {
		field := formatNamespace(e.Namespace())
		errs[field] = translateTag(e)
	}

	return errs
}

func formatNamespace(ns string) string {
	return ns[strings.Index(ns, ".")+1:]
}

func translateTag(fErr validator.FieldError) string {
	switch fErr.Tag() {
	case "required":
		return "obrigatório"

	case "notblank":
		return "não pode estar em branco"

	case "uuid":
		return "deve ser um uuid"

	case "email":
		return "deve ser um e-mail válido"

	case "time":
		return "deve ser um horário válido"

	case "latitude":
		return "deve estar entre -90 e 90"

	case "longitude":
		return "deve estar entre -180 e 180"

	case "numeric":
		return "deve conter apenas dígitos"

	case "alpha":
		return "deve conter apenas letras"

	case "min":
		if fErr.Kind() == reflect.Slice {
			return fmt.Sprintf("deve conter no mínimo %s elemento(s)", fErr.Param())
		}

		return fmt.Sprintf("deve conter no mínimo %s caracteres", fErr.Param())

	case "max":
		if fErr.Kind() == reflect.Slice {
			return fmt.Sprintf("deve conter no máximo %s elemento(s)", fErr.Param())
		}

		return fmt.Sprintf("deve conter no máximo %s caracteres", fErr.Param())

	case "len":
		if fErr.Kind() == reflect.Slice {
			return fmt.Sprintf("deve conter exatamente %s elemento(s)", fErr.Param())
		}

		return fmt.Sprintf("deve conter exatamente %s caractere(s)", fErr.Param())

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
