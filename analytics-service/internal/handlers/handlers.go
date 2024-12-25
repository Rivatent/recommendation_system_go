package handlers

import (
	"analytics-service/internal/model"
	"analytics-service/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// IAnalyticsService - интерфейс для получения аналитических данных.
// Определяет метод GetAnalytics для доступа к данным аналитики.
type IAnalyticsService interface {
	GetAnalytics() ([]model.Analytics, error)
}

// Handler - структура для обработки HTTP-запросов.
// Содержит логгер и сервис аналитики для выполнения операций.
type Handler struct {
	logger log.Factory
	svc    IAnalyticsService
}

// New - функция для создания нового экземпляра обработчика.
// Принимает логгер и сервис аналитики и возвращает указатель на новый Handler.
func New(logger log.Factory, svc IAnalyticsService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}

}

// GetAnalytics - HTTP-обработчик для получения аналитических данных.
// Вызывает метод GetAnalytics у сервиса и возвращает данные в формате JSON.
// В случае ошибки возвращает статус 500 и сообщение об ошибке.
func (h *Handler) GetAnalytics(c *gin.Context) {
	analytics, err := h.svc.GetAnalytics()
	if err != nil {
		h.logger.Bg().Error("failed GetAnalytics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analytics)
}
