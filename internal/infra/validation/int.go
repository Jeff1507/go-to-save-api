package validation

import "fmt"

type IntValidator struct {
	errs  Errors
	field string
	value int
	done  bool
}

func Int(errs Errors, field string, value int) *IntValidator {
	return &IntValidator{
		errs:  errs,
		field: field,
		value: value,
	}
}

func (v *IntValidator) fail(condition bool, message string) *IntValidator {
	if v.done {
		return v
	}

	if condition {
		v.errs.Add(v.field, message)
		v.done = true
	}

	return v
}

func (v *IntValidator) GT(n int) *IntValidator {
	return v.fail(
		v.value <= n,
		fmt.Sprintf("deve ser maior que %v", n),
	)
}

func (v *IntValidator) LT(n int) *IntValidator {
	return v.fail(
		v.value >= n,
		fmt.Sprintf("deve ser menor que %v", n),
	)
}

func (v *IntValidator) GTE(n int) *IntValidator {
	return v.fail(
		v.value < n,
		fmt.Sprintf("deve ser maior ou igual a %v", n),
	)
}

func (v *IntValidator) LTE(n int) *IntValidator {
	return v.fail(
		v.value > n,
		fmt.Sprintf("deve ser menor ou igual a %v", n),
	)
}
