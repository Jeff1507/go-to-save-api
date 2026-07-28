package auth

import (
	"context"
	"errors"
	"time"
)

const (
	sessionTTL       = time.Hour * 24 * 30 // 30 dias
	renewalThreshold = time.Hour * 24 * 15 // 15 dias
)

var ErrInvalidSession = errors.New("invalid or expired session")

type Authenticator interface {
	CreateSession(ctx context.Context, userID string) (string, error)
	ValidateSession(ctx context.Context, sessionID string) (string, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllSessions(ctx context.Context, userID string) error
}
