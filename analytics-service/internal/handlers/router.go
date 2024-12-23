package handlers

import (
	"analytics-service/internal/service"
	"analytics-service/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

func AddHandlers(router *gin.Engine, l log.Factory, svc *service.Service) {
	handlers := New(l, svc)

	app := router.Group("/api/v1")
	{
		userSvc := app.Group("/analytics")
		{
			userSvc.GET("", handlers.GetAnalytics)
		}
	}
}
