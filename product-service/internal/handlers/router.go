package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"product-service/internal/service"
	"product-service/log"
)

func newRouter(l log.Factory, svc *service.Service) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "GET", "PATCH", "DELETE"},
		}),
	)

	AddHandlers(router, l, svc)

	return router
}

func AddHandlers(router *gin.Engine, l log.Factory, svc *service.Service) {
	handlers := New(l, svc)

	app := router.Group("/api/v1")
	{
		ProductSvc := app.Group("/products")
		{
			ProductSvc.GET("", handlers.GetProducts)
			ProductSvc.POST("", handlers.CreateProduct)
			ProductSvc.PATCH("", handlers.UpdateProduct)
			ProductSvc.GET("/:id", handlers.GetProductByID)
			ProductSvc.DELETE("/:id", handlers.DeleteProductByID)
		}
	}
}
