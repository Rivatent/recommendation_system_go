package handlers

import (
	"analytics-service/internal/service"
	"analytics-service/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// newRouter - создает новый маршрутизатор для HTTP-сервера.
// Принимает логгер и сервис аналитики в качестве аргументов.
// Возвращает указатель на созданный экземпляр gin.Engine.
func newRouter(l log.Factory, svc *service.Service) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET"},
		}),
	)

	AddHandlers(router, l, svc)

	return router
}

// AddHandlers - добавляет маршруты и соответствующие обработчики в маршрутизатор.
// Принимает указатель на экземпляр gin.Engine, логгер и сервис аналитики в качестве аргументов.
func AddHandlers(router *gin.Engine, l log.Factory, svc *service.Service) {
	handlers := New(l, svc)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	app := router.Group("/api/v1")
	{
		userSvc := app.Group("/analytics")
		{
			userSvc.GET("", handlers.GetAnalytics)
		}
	}
}
