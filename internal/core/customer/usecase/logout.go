package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
)

type Logout struct {
	sessionManager session.Manager
}

func NewLogout(sessionManager session.Manager) *Logout {
	return &Logout{
		sessionManager: sessionManager,
	}
}

func (uc *Logout) Execute(ctx context.Context, sessionID string) error {
	if err := uc.sessionManager.Revoke(ctx, sessionID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
