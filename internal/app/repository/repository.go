package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db        *pgxpool.Pool
	redisConn *redis.Client
}

func NewRepository(db *pgxpool.Pool, redisConn *redis.Client) *Repository {
	return &Repository{db: db, redisConn: redisConn}
}
