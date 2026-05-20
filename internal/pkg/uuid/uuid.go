package uuid

import "github.com/google/uuid"

func Generate() string {
	return uuid.NewString()
}

func IsValid(value string) bool {
	if err := uuid.Validate(value); err != nil {
		return false
	}

	return true
}
