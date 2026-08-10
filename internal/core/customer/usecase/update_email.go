package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type UpdateEmailInput struct {
	ID    string
	Email string
}

type UpdateEmail struct {
	customerRepo customer.Repository
	otps         otp.Manager
	clock        clock.Clock
}

func NewUpdateEmail(
	customerRepo customer.Repository,
	otps otp.Manager,
	clock clock.Clock,
) *UpdateEmail {
	return &UpdateEmail{
		customerRepo: customerRepo,
		otps:         otps,
		clock:        clock,
	}
}

func (uc *UpdateEmail) Execute(ctx context.Context, input UpdateEmailInput) error {
	now := uc.clock.Now()

	exists, err := uc.customerRepo.Exists(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("check if the email is already in use: %w", err)
	}

	if exists {
		return customer.ErrEmailAlreadyInUse
	}

	cust, err := uc.customerRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find customer by id: %w", err)
	}

	if err := cust.ChangeEmail(input.Email, now); err != nil {
		return fmt.Errorf("change email: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	code, err := uc.otps.Create(ctx, otp.EmailVerification, input.Email)
	if err != nil {
		fmt.Printf("warn: failed to create otp for %s: %v\n", cust.Email(), err)
	} else {
		// TODO: Enviar email com código OTP
		fmt.Printf("Seu código é: %s\n", code)
	}

	return nil
}
