package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
)

type LogoutAll struct {
	sessionManager session.Manager
}

func NewLogoutAll(sessionManager session.Manager) *LogoutAll {
	return &LogoutAll{
		sessionManager: sessionManager,
	}
}

func (uc *LogoutAll) Execute(ctx context.Context, userID string) error {
	if err := uc.sessionManager.RevokeAll(ctx, userID); err != nil {
		return fmt.Errorf("revoke all sessions: %w", err)
	}

	return nil
}
