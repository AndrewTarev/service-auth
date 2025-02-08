package http

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"service-auth/internal/app/delivery/middleware"
	"service-auth/internal/app/service"
	"service-auth/internal/configs"

	docs "service-auth/docs"
)

type Handler struct {
	services service.Service
	cfg      *configs.Config
}

func NewHandler(services service.Service, cfg *configs.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	docs.SwaggerInfo.BasePath = "/auth"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.DELETE("/revoke-token", h.Revoke)
	}
	return router
}
