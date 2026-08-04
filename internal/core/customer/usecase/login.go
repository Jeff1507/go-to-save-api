package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/hash"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	SessionID string
}

type Login struct {
	customerRepo   customer.Repository
	sessionManager session.Manager
}

func NewLogin(customerRepo customer.Repository, sessionManager session.Manager) *Login {
	return &Login{
		customerRepo:   customerRepo,
		sessionManager: sessionManager,
	}
}

func (uc *Login) Execute(ctx context.Context, input LoginInput) (LoginOutput, error) {
	cust, err := uc.customerRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return LoginOutput{}, customer.ErrInvalidCredentials
		}

		return LoginOutput{}, fmt.Errorf("find customer by email: %w", err)
	}

	if !hash.Compare(cust.PasswordHash(), input.Password) {
		return LoginOutput{}, customer.ErrInvalidCredentials
	}

	sessionID, err := uc.sessionManager.Create(ctx, cust.ID())
	if err != nil {
		return LoginOutput{}, fmt.Errorf("create session: %w", err)
	}

	return LoginOutput{
		SessionID: sessionID,
	}, nil
}
