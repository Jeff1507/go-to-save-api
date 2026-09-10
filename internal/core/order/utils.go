package order

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const max = 10_000 // 10^4

func GenerateConfirmationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}

	return fmt.Sprintf("%04d", n.Int64()), nil
}
