package validation

import (
	"encoding/json"
	"strings"
)

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

func (errs Errors) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)

	for f, msg := range errs {
		parts := strings.Split(f, ".")

		current := result // Ponteiro para result

		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = msg
			} else {
				if _, ok := current[part]; !ok {
					current[part] = make(map[string]any)
				}

				current = current[part].(map[string]any)
			}
		}
	}

	return json.Marshal(result)
}
