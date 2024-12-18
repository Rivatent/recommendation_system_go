package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"user-service/internal/service"
)

func newRouter(svc *service.Service) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "GET", "PUT", "DELETE"},
		}),
	)

	AddHandlers(router, svc)

	return router
}

func AddHandlers(router *gin.Engine, svc *service.Service) {
	handlers := New(svc)

	app := router.Group("/api/v1")
	{
		userSvc := app.Group("/users")
		{
			userSvc.GET("", handlers.GetUsers)
			userSvc.POST("", handlers.CreateUser)
			userSvc.PUT("", handlers.UpdateUser)
			userSvc.GET("/:id", handlers.GetUserByID)
			//DELETE?
		}
	}
}
