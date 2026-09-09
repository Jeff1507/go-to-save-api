package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/hash"
)

type ResetPasswordInput struct {
	Email       string
	Code        string
	NewPassword string
}

type ResetPassword struct {
	storeRepo store.Repository
	sessions  session.Manager
	otps      otp.Manager
	clock     clock.Clock
}

func NewResetPassword(
	storeRepo store.Repository,
	sessions session.Manager,
	otps otp.Manager,
	clock clock.Clock,
) *ResetPassword {
	return &ResetPassword{
		storeRepo: storeRepo,
		sessions:  sessions,
		otps:      otps,
		clock:     clock,
	}
}

func (uc *ResetPassword) Execute(ctx context.Context, input ResetPasswordInput) error {
	now := uc.clock.Now()

	if err := uc.otps.Verify(ctx, otp.PasswordReset, input.Email, input.Code); err != nil {
		return fmt.Errorf("verify otp: %w", err)
	}

	s, err := uc.storeRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("find store by email: %w", err)
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
		fmt.Printf("warn: failed to revoke sessions: %v\n", err)
	}

	return nil
}
