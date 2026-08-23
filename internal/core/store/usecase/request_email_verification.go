package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type RequestEmailVerification struct {
	storeRepo store.Repository
	otps      otp.Manager
}

func NewRequestEmailVerification(
	storeRepo store.Repository,
	otps otp.Manager,
) *RequestEmailVerification {
	return &RequestEmailVerification{
		storeRepo: storeRepo,
		otps:      otps,
	}
}

func (uc *RequestEmailVerification) Execute(ctx context.Context, email string) error {
	s, err := uc.storeRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Retorna nil intencionalmente — não revela se o email existe
			return nil
		}

		return fmt.Errorf("find store by email: %w", err)
	}

	if s.IsVerified() {
		// Retorna nil intencionalmente — não revela se o email existe
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
