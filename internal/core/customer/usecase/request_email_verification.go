package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type RequestEmailVerification struct {
	customerRepo customer.Repository
	otps         otp.Manager
}

func NewRequestEmailVerification(
	customerRepo customer.Repository,
	otps otp.Manager,
) *RequestEmailVerification {
	return &RequestEmailVerification{
		customerRepo: customerRepo,
		otps:         otps,
	}
}

func (uc *RequestEmailVerification) Execute(ctx context.Context, email string) error {
	cust, err := uc.customerRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			// Retorna nil intencionalmente — não revela se o email existe
			return nil
		}

		return fmt.Errorf("find customer by email: %w", err)
	}

	if cust.IsVerified() {
		// Retorna nil intencionalmente — não revela se já está verificado
		return nil
	}

	code, err := uc.otps.Create(ctx, otp.EmailVerification, email)
	if err != nil {
		return fmt.Errorf("create otp: %w", err)
	}

	// TODO: Enviar email com código OTP
	fmt.Printf("Seu código é: %s\n", code)

	return nil
}
