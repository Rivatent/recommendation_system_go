package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"recommendation-service/internal/service"
	"recommendation-service/pkg/log"
)

// newRouter - функция, создающая новый маршрутизатор для приложения.
// Принимает логгер и сервис в качестве параметров и возвращает созданный
// экземпляр *gin.Engine, который представляет собой маршрутизатор Gin.
// l - фабрика логирования для ведения записи о событиях и ошибках.
// svc - указатель на сервис, который будет использоваться для обработки логики приложения.
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

// AddHandlers - функция, отвечающая за добавление различных обработчиков
// для API-маршрутов в маршрутизатор.
// Принимает маршрутизатор, логгер и сервис в качестве аргументов.
// router - экземпляр маршрутизатора Gin, к которому добавляются маршруты.
// l - фабрика логирования для ведения записей о событиях и ошибках.
// svc - указатель на сервис, который будет использоваться для обработки логики приложения.
func AddHandlers(router *gin.Engine, l log.Factory, svc *service.Service) {
	handlers := New(l, svc)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	app := router.Group("/api/v1")
	{
		userSvc := app.Group("/recommendations")
		{
			userSvc.GET("", handlers.GetRecommendations)
			userSvc.GET("/users/:id", handlers.GetRecommendationsByUserID)
			userSvc.GET("/recs/:id", handlers.GetRecommendationByID)
		}
	}
}
