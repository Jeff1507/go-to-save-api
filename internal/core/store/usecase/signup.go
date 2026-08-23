package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/hash"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type SignupInput struct {
	Email    string
	Name     string
	CNPJ     string
	Password string
	Address  AddressInput
}

type SignupOutput struct {
	UserID    string
	SessionID string
}

type Signup struct {
	storeRepo store.Repository
	sessions  session.Manager
	otps      otp.Manager
	clock     clock.Clock
}

func NewSignup(
	storeRepo store.Repository,
	sessions session.Manager,
	otps otp.Manager,
	clock clock.Clock,
) *Signup {
	return &Signup{
		storeRepo: storeRepo,
		sessions:  sessions,
		otps:      otps,
		clock:     clock,
	}
}

func (uc *Signup) Execute(ctx context.Context, input SignupInput) (SignupOutput, error) {
	now := uc.clock.Now()

	exists, err := uc.storeRepo.Exists(ctx, input.Email)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("check if the email is already in use: %w", err)
	}

	if exists {
		return SignupOutput{}, store.ErrEmailAlreadyInUse
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("hash password: %w", err)
	}

	address, err := store.NewAddress(
		input.Address.PostalCode,
		input.Address.Street,
		input.Address.Number,
		input.Address.Neighborhood,
		input.Address.Complement,
		input.Address.City,
		input.Address.State,
		input.Address.Latitude,
		input.Address.Longitude,
	)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("new address: %w", err)
	}

	s, err := store.New(
		uuid.Generate(),
		input.Email,
		input.Name,
		input.CNPJ,
		passwordHash,
		address,
		now,
	)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("new store: %w", err)
	}

	if err := uc.storeRepo.Save(ctx, s); err != nil {
		return SignupOutput{}, fmt.Errorf("save store: %w", err)
	}

	sessionID, err := uc.sessions.Create(ctx, s.ID())
	if err != nil {
		return SignupOutput{}, fmt.Errorf("create session: %w", err)
	}

	// Falha silenciosa, usuário pode pedir reenvio depois
	code, err := uc.otps.Create(ctx, otp.EmailVerification, s.Email())
	if err != nil {
		fmt.Printf("warn: failed to create otp for %s: %v\n", s.Email(), err)
	} else {
		// TODO: Enviar email com código OTP
		fmt.Printf("Seu código é: %s\n", code)
	}

	return SignupOutput{
		UserID:    s.ID(),
		SessionID: sessionID,
	}, nil
}
