package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/hash"
)

type UpdatePasswordInput struct {
	ID          string
	Password    string
	NewPassword string
}

type UpdatePassword struct {
	customerRepo customer.Repository
	sessions     session.Manager
	clock        clock.Clock
}

func NewUpdatePassword(
	customerRepo customer.Repository,
	sessions session.Manager,
	clock clock.Clock,
) *UpdatePassword {
	return &UpdatePassword{
		customerRepo: customerRepo,
		sessions:     sessions,
		clock:        clock,
	}
}

func (uc *UpdatePassword) Execute(ctx context.Context, input UpdatePasswordInput) error {
	now := uc.clock.Now()

	cust, err := uc.customerRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find customer by id: %w", err)
	}

	passwordsMatch := hash.Compare(cust.PasswordHash(), input.Password)
	if !passwordsMatch {
		return customer.ErrIncorrectPassword
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
