package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"product-service/internal/service"
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
		userSvc := app.Group("/products")
		{
			userSvc.GET("", handlers.GetProducts)
			//userSvc.POST("", handlers.AddProduct)
			//userSvc.PUT("", handlers.UpdateProduct)
			//userSvc.DELETE("/:id", handlers.DeleteProduct)
		}
	}
}
