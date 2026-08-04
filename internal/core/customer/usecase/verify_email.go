package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type VerifyEmailInput struct {
	Email string
	Code  string
}

type VerifyEmail struct {
	clock        clock.Clock
	customerRepo customer.Repository
	otpManager   otp.Manager
}

func (uc *VerifyEmail) Execute(ctx context.Context, input VerifyEmailInput) error {
	now := uc.clock.Now()

	if err := uc.otpManager.Verify(
		ctx,
		otp.EmailVerification,
		input.Email,
		input.Code,
	); err != nil {
		return fmt.Errorf("verify otp: %w", err)
	}

	cust, err := uc.customerRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("find customer by email: %w", err)
	}

	if err := cust.Verify(now); err != nil {
		return fmt.Errorf("verify email: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	return nil
}
