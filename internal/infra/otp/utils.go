package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const max = 1_000_000 // 10^6

type Purpose uint8

const (
	Unknown Purpose = iota
	EmailVerification
	PasswordReset
)

func (p Purpose) String() string {
	switch p {
	case EmailVerification:
		return "email_verification"
	case PasswordReset:
		return "password_reset"
	default:
		return "unknown"
	}
}

func ParsePurpose(str string) (Purpose, error) {
	switch str {
	case "email_verification":
		return EmailVerification, nil
	case "password_reset":
		return PasswordReset, nil
	default:
		return Unknown, fmt.Errorf("unknown purpose: %q", str)
	}
}

func GenerateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return "", fmt.Errorf("generate otp: %v", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
