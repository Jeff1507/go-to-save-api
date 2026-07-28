package otp

import (
	"context"
	"errors"
	"time"
)

const TTL = time.Minute * 15

var ErrInvalidOTP = errors.New("invalid or expired otp")

type Manager interface {
	Create(ctx context.Context, purpose Purpose, email string) (string, error)
	Verify(ctx context.Context, purpose Purpose, email, code string) error
}
