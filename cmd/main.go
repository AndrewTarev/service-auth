package main

import (
	"fmt"

	logger "github.com/sirupsen/logrus"

	"service-auth/internal/app/delivery/http"
	"service-auth/internal/app/repository"
	"service-auth/internal/app/service"
	"service-auth/internal/app/utils"
	"service-auth/internal/configs"
	"service-auth/internal/server"
	"service-auth/pkg/db"
	logging "service-auth/pkg/logger"
	"service-auth/pkg/redis_client"
)

// @title           Auth
// @version         1.0
// @description     API для авторизации и аутентификации.

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig("./internal/configs")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	// Настройка логгера
	logging.SetupLogger(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.OutputFile)

	// Подключение к базе данных
	dbConn, err := db.ConnectPostgres(cfg.Database.Dsn)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close()

	redisConn, err := redis_client.InitRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}

	// db.ApplyMigrations(cfg.Database.Dsn, cfg.Database.MigratePath)

	// загрузка auth параметров
	jwtManager := utils.NewJWTManager(cfg)
	repo := repository.NewRepository(dbConn, redisConn)
	services := service.NewService(repo, jwtManager, cfg)
	handlers := http.NewHandler(services, cfg)

	// Настройка и запуск сервера
	server.SetupAndRunServer(&cfg.Server, handlers.InitRoutes())
}
