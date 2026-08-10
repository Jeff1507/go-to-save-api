package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
)

type LogoutAll struct {
	sessions session.Manager
}

func NewLogoutAll(sessions session.Manager) *LogoutAll {
	return &LogoutAll{
		sessions: sessions,
	}
}

func (uc *LogoutAll) Execute(ctx context.Context, userID string) error {
	if err := uc.sessions.RevokeAll(ctx, userID); err != nil {
		return fmt.Errorf("revoke all sessions: %w", err)
	}

	return nil
}
