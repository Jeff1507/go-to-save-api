package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
)

type Logout struct {
	sessions session.Manager
}

func NewLogout(sessions session.Manager) *Logout {
	return &Logout{
		sessions: sessions,
	}
}

func (uc *Logout) Execute(ctx context.Context, sessionID string) error {
	if err := uc.sessions.Revoke(ctx, sessionID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
