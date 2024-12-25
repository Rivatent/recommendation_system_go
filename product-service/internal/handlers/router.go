package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"product-service/internal/service"
	"product-service/log"
)

// newRouter создает новый экземпляр Gin-роутера, настраивает обработчики
// промежуточного ПО и добавляет маршруты для API.
// Параметры:
//   - l: экземпляр логгера, используемый для записи логов.
//   - svc: указатель на сервис, содержащий бизнес-логику для работы с продуктами.
//
// Возвращает указатель на настроенный экземпляр gin.Engine.
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

// AddHandlers настраивает маршруты и назначает обработчики для операций
// с продуктами в приложении. Создает группы маршрутов для API и
// связывает их с соответствующими методами обработчиков.
// Параметры:
//   - router: экземпляр gin.Engine, на который добавляются маршруты.
//   - l: экземпляр логгера для записи логов.
//   - svc: указатель на сервис, предоставляющий бизнес-логику для работы с продуктами.
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
