package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisAuthenticator struct {
	rdb *redis.Client
}

func NewRedisAuthenticator(rdb *redis.Client) *RedisAuthenticator {
	return &RedisAuthenticator{
		rdb: rdb,
	}
}

func (a RedisAuthenticator) key(sessionID string) string {
	return "session" + ":" + sessionID
}

func (a RedisAuthenticator) sessionsKey(userID string) string {
	return "user_sessions" + ":" + userID
}

func (a *RedisAuthenticator) CreateSession(ctx context.Context, userID string) (string, error) {
	sessionID := uuid.Generate()

	key := a.key(sessionID)
	sessionsKey := a.sessionsKey(userID)

	tx := a.rdb.TxPipeline()

	tx.Set(ctx, key, userID, sessionTTL)

	tx.SAdd(ctx, sessionsKey, key)

	if _, err := tx.Exec(ctx); err != nil {
		return "", fmt.Errorf("execute pipeline: %w", err)
	}

	return sessionID, nil
}

func (a *RedisAuthenticator) ValidateSession(
	ctx context.Context,
	sessionID string,
) (string, error) {
	key := a.key(sessionID)

	userID, err := a.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrInvalidSession
		}

		return "", fmt.Errorf("get session: %w", err)
	}

	if ttl := a.rdb.TTL(ctx, key).Val(); ttl > 0 && ttl <= renewalThreshold {
		a.rdb.Set(ctx, key, userID, sessionTTL) // Falha silenciosa
	}

	return userID, nil
}

func (a *RedisAuthenticator) RevokeSession(ctx context.Context, sessionID string) error {
	key := a.key(sessionID)

	if err := a.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (a *RedisAuthenticator) RevokeAllSessions(ctx context.Context, userID string) error {
	sessionsKey := a.sessionsKey(userID)

	sessions, err := a.rdb.SMembers(ctx, sessionsKey).Result()
	if err != nil {
		return fmt.Errorf("get user sessions: %w", err)
	}

	if len(sessions) == 0 {
		return nil
	}

	tx := a.rdb.TxPipeline()

	for _, session := range sessions {
		tx.Del(ctx, session)
	}

	tx.Del(ctx, sessionsKey)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("delete user sessions: %w", err)
	}

	return nil
}
