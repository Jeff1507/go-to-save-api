package otp

import (
	"context"
	"errors"
	"fmt"

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

func (mng RedisManager) otpKey(purpose Purpose, email string) string {
	return "otp" + ":" + purpose.String() + ":" + email // "otp:purpose:email"
}

func (mng *RedisManager) Create(
	ctx context.Context,
	purpose Purpose,
	email string,
) (string, error) {
	otpKey := mng.otpKey(purpose, email)

	code, err := GenerateCode()
	if err != nil {
		return "", fmt.Errorf("generate random code: %w", err)
	}

	if err := mng.rdb.Set(ctx, otpKey, code, TTL).Err(); err != nil {
		return "", fmt.Errorf("insert otp: %w", err)
	}

	return code, nil
}

func (mng *RedisManager) Verify(ctx context.Context, purpose Purpose, email, code string) error {
	key := mng.otpKey(purpose, email)

	otpCode, err := mng.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrInvalidOTP
		}

		return fmt.Errorf("get otp: %w", err)
	}

	if otpCode != code {
		return ErrInvalidOTP
	}

	mng.rdb.Del(ctx, key)

	return nil
}
