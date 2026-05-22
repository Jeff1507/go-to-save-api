package validation

import (
	"fmt"
	"strings"

	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type StringValidator struct {
	errs  Errors
	field string
	value string
	done  bool
}

func String(errs Errors, field, value string) *StringValidator {
	return &StringValidator{
		errs:  errs,
		field: field,
		value: value,
	}
}

func (v *StringValidator) fail(condition bool, message string) *StringValidator {
	if v.done {
		return v
	}

	if condition {
		v.errs.Add(v.field, message)
		v.done = true
	}

	return v
}

func isEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

func (v *StringValidator) Optional() *StringValidator {
	if isEmpty(v.value) {
		v.done = true
	}

	return v
}

func (v *StringValidator) Required() *StringValidator {
	return v.fail(isEmpty(v.value), "obrigatório")
}

func (v *StringValidator) UUID() *StringValidator {
	return v.fail(!uuid.IsValid(v.value), "deve ser um uuid")
}

func (v *StringValidator) Min(n int) *StringValidator {
	return v.fail(
		len([]rune(v.value)) < n,
		fmt.Sprintf("deve ter no mínimo %d caracteres", n),
	)
}

func (v *StringValidator) Max(n int) *StringValidator {
	return v.fail(
		len([]rune(v.value)) > n,
		fmt.Sprintf("deve ter no máximo %d caracteres", n),
	)
}
