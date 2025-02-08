package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"

	"service-auth/internal/app/errs"
)

const ActiveToken = 1

type RedisRepo struct {
	redisConn *redis.Client
}

func NewRedisRepo(redisConn *redis.Client) *RedisRepo {
	return &RedisRepo{redisConn: redisConn}
}

func (r *RedisRepo) SaveToken(ctx context.Context, token string, ttl time.Duration) error {
	err := r.redisConn.Set(ctx, token, ActiveToken, ttl).Err()
	if err != nil {
		logger.Errorf("save token error: %v", err)
		return errs.ErrFailedToSave
	}
	logger.Debugf("save token: %v", token)
	return nil
}

func (r *RedisRepo) FindTokenInRedis(ctx context.Context, token string) error {
	logger.Debugf("Checking token %s", token)
	_, err := r.redisConn.Get(ctx, token).Int()
	if errors.Is(err, redis.Nil) {
		logger.Debugf("Token %s does not exist in Redis", token)
		return errs.ErrTokenNotFound
	} else if err != nil {
		logger.Errorf("Failed to get token from Redis: %w", err)
		return errs.ErrValidateInRedis
	}
	return nil
}

func (r *RedisRepo) DeleteToken(ctx context.Context, token string) error {
	deleted, err := r.redisConn.Del(ctx, token).Result()
	if err != nil {
		logger.Errorf("Failed to delete token from Redis (token: %s): %v", token, err)
		return fmt.Errorf("failed to delete token: %w", err)
	}
	logger.Debugf("Deleted token from Redis (token: %s)", token)
	if deleted == 0 {
		logger.Warnf("Token not found in Redis: %s", token)
		return errs.ErrTokenNotFound
	}

	logger.Debugf("Deleted token from Redis: %s", token)
	return nil
}

func (r *RedisRepo) UpdateRefreshTokenInRedis(ctx context.Context, oldToken, newToken string, ttl time.Duration) error {
	// Начало транзакции
	pipe := r.redisConn.TxPipeline()

	// Удаление старого токена
	pipe.Del(ctx, oldToken)

	// Сохранение нового токена
	pipe.Set(ctx, newToken, ActiveToken, ttl)

	// Выполнение транзакции
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Errorf("Failed to update token in Redis: %s", err)
		return errs.ErrFailedToRefresh
	}
	logger.Debugf("Successfully updated token in Redis")
	return nil
}
