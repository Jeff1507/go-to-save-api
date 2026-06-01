package validation

import (
	"fmt"
	"time"
)

type TimeValidator struct {
	errs  Errors
	field string
	value time.Time
	done  bool
}

func Time(errs Errors, field string, value time.Time) *TimeValidator {
	return &TimeValidator{
		errs:  errs,
		field: field,
		value: value,
	}
}

func (v *TimeValidator) fail(condition bool, message string) *TimeValidator {
	if v.done {
		return v
	}

	if condition {
		v.errs.Add(v.field, message)
		v.done = true
	}

	return v
}

func (v *TimeValidator) Zero() *TimeValidator {
	if v.value.IsZero() {
		v.done = true
	}

	return v
}

func (v *TimeValidator) NotZero() *TimeValidator {
	return v.fail(v.value.IsZero(), "obrigatório")
}

func (v *TimeValidator) Before(t time.Time) *TimeValidator {
	var msg string

	if t.IsZero() {
		msg = "inválido"
	} else {
		msg = fmt.Sprintf("deve ser anterior a %s", t.Format("02/01/2006 15:04"))
	}

	return v.fail(!v.value.Before(t), msg)
}

func (v *TimeValidator) After(t time.Time) *TimeValidator {
	var msg string

	if t.IsZero() {
		msg = "inválido"
	} else {
		msg = fmt.Sprintf("deve vir depois de %s", t.Format("02/01/2006 15:04"))
	}

	return v.fail(!v.value.After(t), msg)
}
