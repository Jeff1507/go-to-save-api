package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type RequestPasswordReset struct {
	customerRepo customer.Repository
	otpManager   otp.Manager
}

func NewRequestPasswordReset(
	customerRepo customer.Repository,
	otpManager otp.Manager,
) *RequestPasswordReset {
	return &RequestPasswordReset{
		customerRepo: customerRepo,
		otpManager:   otpManager,
	}
}

func (uc *RequestPasswordReset) Execute(ctx context.Context, email string) error {
	exists, err := uc.customerRepo.Exists(ctx, email)
	if err != nil {
		return fmt.Errorf("check if the email exists: %w", err)
	}

	if !exists {
		// Retorna nil intencionalmente — não revela se o email existe
		return nil
	}

	code, err := uc.otpManager.Create(ctx, otp.PasswordReset, email)
	if err != nil {
		return fmt.Errorf("create otp: %w", err)
	}

	// TODO: Enviar email com código OTP
	fmt.Printf("Seu código é: %s\n", code)

	return nil
}
