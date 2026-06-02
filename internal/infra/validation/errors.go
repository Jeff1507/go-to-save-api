package validation

import "strings"

type Errors map[string]string

func NewErrors() Errors {
	return make(Errors)
}

func (errs Errors) Error() string {
	var msgs []string

	for f, msg := range errs {
		msgs = append(msgs, f+": "+msg)
	}

	return strings.Join(msgs, "; ")
}

func (errs Errors) Add(field, message string) {
	errs[field] = message
}

func (errs Errors) HasErrors() bool {
	return len(errs) > 0
}

func (errs Errors) Has(field string) bool {
	_, ok := errs[field]
	return ok
}
