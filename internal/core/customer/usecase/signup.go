package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/hash"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type SignupInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type SignupOutput struct {
	UserID    string
	SessionID string
}

type Signup struct {
	customerRepo customer.Repository
	sessions     session.Manager
	otps         otp.Manager
	clock        clock.Clock
}

func NewSignup(
	customerRepo customer.Repository,
	sessions session.Manager,
	otps otp.Manager,
	clock clock.Clock,
) *Signup {
	return &Signup{
		customerRepo: customerRepo,
		sessions:     sessions,
		otps:         otps,
		clock:        clock,
	}
}

func (uc *Signup) Execute(ctx context.Context, input SignupInput) (SignupOutput, error) {
	now := uc.clock.Now()

	exists, err := uc.customerRepo.Exists(ctx, input.Email)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("check if the email is already in use: %w", err)
	}

	if exists {
		return SignupOutput{}, customer.ErrEmailAlreadyInUse
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("hash password: %w", err)
	}

	cust, err := customer.New(
		uuid.Generate(),
		input.Email,
		input.FirstName,
		input.LastName,
		passwordHash,
		now,
	)
	if err != nil {
		return SignupOutput{}, fmt.Errorf("new customer: %w", err)
	}

	if err := uc.customerRepo.Save(ctx, cust); err != nil {
		return SignupOutput{}, fmt.Errorf("save customer: %w", err)
	}

	sessionID, err := uc.sessions.Create(ctx, cust.ID())
	if err != nil {
		return SignupOutput{}, fmt.Errorf("create session: %w", err)
	}

	// Falha silenciosa, usuário pode pedir reenvio depois
	code, err := uc.otps.Create(ctx, otp.EmailVerification, cust.Email())
	if err != nil {
		fmt.Printf("warn: failed to create otp for %s: %v\n", cust.Email(), err)
	} else {
		// TODO: Enviar email com código OTP
		fmt.Printf("Seu código é: %s\n", code)
	}

	return SignupOutput{
		UserID:    cust.ID(),
		SessionID: sessionID,
	}, nil
}
