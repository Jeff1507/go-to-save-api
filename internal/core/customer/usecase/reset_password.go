package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
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
	customerRepo customer.Repository
	sessions     session.Manager
	otps         otp.Manager
	clock        clock.Clock
}

func NewResetPassword(
	customerRepo customer.Repository,
	sessions session.Manager,
	otps otp.Manager,
	clock clock.Clock,
) *ResetPassword {
	return &ResetPassword{
		customerRepo: customerRepo,
		sessions:     sessions,
		otps:         otps,
		clock:        clock,
	}
}

func (uc *ResetPassword) Execute(ctx context.Context, input ResetPasswordInput) error {
	now := uc.clock.Now()

	if err := uc.otps.Verify(ctx, otp.PasswordReset, input.Email, input.Code); err != nil {
		return fmt.Errorf("verify otp: %w", err)
	}

	cust, err := uc.customerRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("find customer by email: %w", err)
	}

	passwordHash, err := hash.Generate(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := cust.ChangePassword(passwordHash, now); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	if err := uc.sessions.RevokeAll(ctx, cust.ID()); err != nil {
		return fmt.Errorf("revoke all sessions: %w", err)
	}

	return nil
}
