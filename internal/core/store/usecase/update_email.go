package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type UpdateEmailInput struct {
	ID    string
	Email string
}

type UpdateEmail struct {
	storeRepo store.Repository
	otps      otp.Manager
	clock     clock.Clock
}

func NewUpdateEmail(storeRepo store.Repository, otps otp.Manager, clock clock.Clock) *UpdateEmail {
	return &UpdateEmail{
		storeRepo: storeRepo,
		otps:      otps,
		clock:     clock,
	}
}

// Atenção: se o usuário tentar alterar para o próprio email atual
// receberá um erro (ErrEmailAlreadyInUse)
func (uc *UpdateEmail) Execute(ctx context.Context, input UpdateEmailInput) error {
	now := uc.clock.Now()

	exists, err := uc.storeRepo.Exists(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("check if the email is already in use: %w", err)
	}

	if exists {
		return store.ErrEmailAlreadyInUse
	}

	s, err := uc.storeRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find store by id: %w", err)
	}

	if err := s.ChangeEmail(input.Email, now); err != nil {
		return fmt.Errorf("change email: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	code, err := uc.otps.Create(ctx, otp.EmailVerification, input.Email)
	if err != nil {
		fmt.Printf("warn: failed to create otp for %s: %v\n", s.Email(), err)
	} else {
		// TODO: Enviar email com código OTP
		fmt.Printf("Seu código é: %s\n", code)
	}

	return nil
}
