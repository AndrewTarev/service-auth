package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	logger "github.com/sirupsen/logrus"

	"service-auth/internal/app/delivery/http"
	"service-auth/internal/app/repository"
	"service-auth/internal/app/service"
	"service-auth/internal/configs"
	"service-auth/internal/infrastructure/auth_helper"
	"service-auth/internal/infrastructure/db"
	logging "service-auth/internal/infrastructure/logger"
	"service-auth/internal/infrastructure/redis_client"
	"service-auth/internal/server"
)

// @title           Auth
// @version         1.0
// @description     API для авторизации и аутентификации.

// @host      localhost:8080
// @BasePath  /auth
func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig("./internal/configs")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	// Настройка логгера
	logging.SetupLogger(&cfg.Logging)

	// Подключение к базе данных
	dbConn, err := db.ConnectPostgres(&cfg.Database)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close()

	redisConn, err := redis_client.InitRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}

	// applyMigrations(&cfg.Database)

	// загрузка auth параметров
	jwtManager := auth_helper.NewJWTManager(cfg)
	repo := repository.NewRepository(dbConn, redisConn)
	services := service.NewService(repo, jwtManager, cfg)
	handlers := http.NewHandler(services, cfg)

	// Настройка и запуск сервера
	server.SetupAndRunServer(&cfg.Server, handlers.InitRoutes())
}

func applyMigrations(cfg *configs.PostgresConfig) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)
	m, err := migrate.New(
		"file:///app/internal/infrastructure/db/migrations",
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
