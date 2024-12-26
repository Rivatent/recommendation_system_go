package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"user-service/internal/service"
	"user-service/pkg/log"
)

// newRouter создает и настраивает новый маршрутизатор Gin.
//func newRouter(l log.Factory, svc *service.Service) *gin.Engine {
//	router := gin.New()
//
//	router.Use(
//		gin.Recovery(),
//		gin.Logger(),
//		cors.New(cors.Config{
//			AllowAllOrigins: true,
//			AllowMethods:    []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS"},
//		}),
//	)
//
//	AddHandlers(router, l, svc)
//
//	return router
//}

// newRouter создает и настраивает новый маршрутизатор Gin.
func newRouter(l log.Factory, svc *service.Service) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
		}),
	)

	AddHandlers(router, l, svc)

	router.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Status(http.StatusNoContent)
	})

	return router
}

// AddHandlers определяет маршруты API для обработки запросов в маршрутизаторе Gin.
func AddHandlers(router *gin.Engine, l log.Factory, svc *service.Service) {
	handlers := New(l, svc)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	app := router.Group("/api/v1")
	{
		userSvc := app.Group("/users")
		{
			userSvc.GET("", handlers.GetUsers)
			userSvc.POST("", handlers.CreateUser)
			userSvc.PATCH("", handlers.UpdateUser)
			userSvc.GET("/:id", handlers.GetUserByID)
		}
	}
}
