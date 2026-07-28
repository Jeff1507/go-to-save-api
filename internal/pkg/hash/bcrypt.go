package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Generate(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash: %v", err)
	}

	return string(passwordHash), nil
}

func Compare(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
