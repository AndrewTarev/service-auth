package db

import (
	"context"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	logger "github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ConnectPostgres(dsn string) (*pgxpool.Pool, error) {
	// Настраиваем контекст для подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем пул соединений
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Errorf("Error parsing PostgreSQL DSN: %v", err)
		return nil, err
	}

	// Применяем дополнительные настройки пула
	poolConfig.MaxConns = 50                       // Максимальное количество соединений
	poolConfig.MinConns = 5                        // Минимальное количество соединений
	poolConfig.HealthCheckPeriod = 1 * time.Minute // Период проверки соединений

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Errorf("Failed to connect to PostgreSQL: %v", err)
		return nil, err
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		logger.Errorf("PostgreSQL ping failed: %v", err)
		return nil, err
	}

	logger.Debug("Successfully connected to PostgreSQL")
	return pool, nil
}

func ApplyMigrations(dsn string, migratePath string) {
	m, err := migrate.New(
		migratePath,
		dsn,
	)
	if err != nil {
		logger.Fatalf("Could not initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatalf("Could not apply migrations: %v", err)
	}

	logger.Debug("Migrations applied successfully!")
}
