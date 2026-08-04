package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	rdb *redis.Client
}

func NewRedisManager(rdb *redis.Client) *RedisManager {
	return &RedisManager{
		rdb: rdb,
	}
}

func (mng RedisManager) sessionKey(id string) string {
	return "session" + ":" + id
}

func (mng RedisManager) userSessionsKey(userID string) string {
	return "user_sessions" + ":" + userID
}

func (mng *RedisManager) Create(ctx context.Context, userID string) (string, error) {
	sessionID := uuid.Generate()

	sessionKey := mng.sessionKey(sessionID)
	userSessionsKey := mng.userSessionsKey(userID)

	tx := mng.rdb.TxPipeline()

	tx.Set(ctx, sessionKey, userID, TTL)

	tx.SAdd(ctx, userSessionsKey, sessionKey)

	if _, err := tx.Exec(ctx); err != nil {
		return "", fmt.Errorf("execute pipeline: %w", err)
	}

	return sessionID, nil
}

func (mng *RedisManager) Validate(ctx context.Context, id string) (string, error) {
	sessionKey := mng.sessionKey(id)

	userID, err := mng.rdb.Get(ctx, sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrInvalidSession
		}

		return "", fmt.Errorf("get session: %w", err)
	}

	if ttl := mng.rdb.TTL(ctx, sessionKey).Val(); ttl > 0 && ttl <= renewalThreshold {
		// Falha silenciosa
		mng.rdb.Set(ctx, sessionKey, userID, TTL)
	}

	return userID, nil
}

func (mng *RedisManager) Revoke(ctx context.Context, id string) error {
	sessionKey := mng.sessionKey(id)

	if err := mng.rdb.Del(ctx, sessionKey).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (mng *RedisManager) RevokeAll(ctx context.Context, userID string) error {
	userSessionsKey := mng.userSessionsKey(userID)

	sessions, err := mng.rdb.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return fmt.Errorf("get user sessions: %w", err)
	}

	if len(sessions) == 0 {
		return nil
	}

	tx := mng.rdb.TxPipeline()

	for _, session := range sessions {
		tx.Del(ctx, session)
	}

	tx.Del(ctx, userSessionsKey)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("execute pipeline: %w", err)
	}

	return nil
}
