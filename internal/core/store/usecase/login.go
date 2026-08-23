package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
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
	storeRepo store.Repository
	sessions  session.Manager
}

func NewLogin(storeRepo store.Repository, sessions session.Manager) *Login {
	return &Login{
		storeRepo: storeRepo,
		sessions:  sessions,
	}
}

func (uc *Login) Execute(ctx context.Context, input LoginInput) (LoginOutput, error) {
	s, err := uc.storeRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return LoginOutput{}, store.ErrInvalidCredentials
		}

		return LoginOutput{}, fmt.Errorf("find store by email: %w", err)
	}

	passwordsMatch := hash.Compare(s.PasswordHash(), input.Password)
	if !passwordsMatch {
		return LoginOutput{}, store.ErrInvalidCredentials
	}

	sessionID, err := uc.sessions.Create(ctx, s.ID())
	if err != nil {
		return LoginOutput{}, fmt.Errorf("create session: %w", err)
	}

	return LoginOutput{
		SessionID: sessionID,
	}, nil
}
