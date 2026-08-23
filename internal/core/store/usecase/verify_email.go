package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type VerifyEmailInput struct {
	Email string
	Code  string
}

type VerifyEmail struct {
	storeRepo store.Repository
	otps      otp.Manager
	clock     clock.Clock
}

func NewVerifyEmail(storeRepo store.Repository, otps otp.Manager, clock clock.Clock) *VerifyEmail {
	return &VerifyEmail{
		storeRepo: storeRepo,
		otps:      otps,
		clock:     clock,
	}
}

func (uc *VerifyEmail) Execute(ctx context.Context, input VerifyEmailInput) error {
	now := uc.clock.Now()

	if err := uc.otps.Verify(ctx, otp.EmailVerification, input.Email, input.Code); err != nil {
		return fmt.Errorf("verify otp: %w", err)
	}

	s, err := uc.storeRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("find store by email: %w", err)
	}

	if err := s.Verify(now); err != nil {
		return fmt.Errorf("verify email: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	return nil
}
