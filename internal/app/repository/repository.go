package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"service-auth/internal/app/models"
)

type PostgresRepository interface {
	CreateUser(ctx context.Context, user models.UserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, username string) (models.GetUserResponse, error)
}

type RedisRepository interface {
	SaveToken(ctx context.Context, token string, ttl time.Duration) error
	FindTokenInRedis(ctx context.Context, token string) error
	DeleteToken(ctx context.Context, token string) error
	UpdateRefreshTokenInRedis(ctx context.Context, oldToken, newToken string, ttl time.Duration) error
}

type Repository struct {
	PostgresRepository
	RedisRepository
}

func NewRepository(db *pgxpool.Pool, redisConn *redis.Client) *Repository {
	return &Repository{
		PostgresRepository: NewPostgresRepo(db),
		RedisRepository:    NewRedisRepo(redisConn),
	}
}
