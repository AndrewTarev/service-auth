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

type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	Revoke(ctx *gin.Context)
}

type Handler struct {
	AuthHandler
}

func NewHandler(services service.Service, cfg *configs.Config) *Handler {
	return &Handler{
		AuthHandler: NewAuth(services, cfg),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.Refresh)
			auth.DELETE("/revoke-token", h.Revoke)
		}
	}

	return router
}
