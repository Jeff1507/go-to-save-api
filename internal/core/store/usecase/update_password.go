package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/hash"
)

type UpdatePasswordInput struct {
	ID          string
	Password    string
	NewPassword string
}

type UpdatePassword struct {
	storeRepo store.Repository
	sessions  session.Manager
	clock     clock.Clock
}

func NewUpdatePassword(
	storeRepo store.Repository,
	sessions session.Manager,
	clock clock.Clock,
) *UpdatePassword {
	return &UpdatePassword{
		storeRepo: storeRepo,
		sessions:  sessions,
		clock:     clock,
	}
}

func (uc *UpdatePassword) Execute(ctx context.Context, input UpdatePasswordInput) error {
	now := uc.clock.Now()

	s, err := uc.storeRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find store by id: %w", err)
	}

	passwordsMatch := hash.Compare(s.PasswordHash(), input.Password)
	if !passwordsMatch {
		return store.ErrIncorrectPassword
	}

	passwordHash, err := hash.Generate(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.ChangePassword(passwordHash, now); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	if err := uc.sessions.RevokeAll(ctx, s.ID()); err != nil {
		return fmt.Errorf("revoke all sessions: %w", err)
	}

	return nil
}
