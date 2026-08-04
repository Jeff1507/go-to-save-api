package session

import (
	"context"
	"errors"
	"time"
)

const (
	TTL              = time.Hour * 24 * 30 // 30 dias
	renewalThreshold = time.Hour * 24 * 15 // 15 dias
)

var ErrInvalidSession = errors.New("invalid or expired session")

type Manager interface {
	Create(ctx context.Context, userID string) (string, error)
	// TODO: Sessões de customer e store compartilham o mesmo Redis sem distinção de role.
	// Uma sessão de store pode ser usada em rotas de customer e vice-versa.
	// Solução: adicionar role no valor da sessão ("customer:userID" / "store:userID")
	// e validar a role no middleware antes de injetar o ID no contexto.
	Validate(ctx context.Context, id string) (string, error)
	Revoke(ctx context.Context, id string) error
	RevokeAll(ctx context.Context, userID string) error
}
