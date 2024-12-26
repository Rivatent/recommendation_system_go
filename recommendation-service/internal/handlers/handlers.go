package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"recommendation-service/internal/model"
	"recommendation-service/internal/monitoring"
	"recommendation-service/log"
	"time"
)

// IRecommendationsService - интерфейс, определяющий методы для работы с рекомендациями.
// Включает методы для получения всех рекомендаций, получения рекомендаций по идентификатору
// и получения рекомендаций по идентификатору пользователя.
type IRecommendationsService interface {
	GetRecommendations() ([]model.Recommendation, error)
	GetRecommendationByID(id string) (model.Recommendation, error)
	GetRecommendationsByUserID(id string) ([]model.Recommendation, error)
}

// Handler - структура, представляющая обработчик HTTP-запросов для получения рекомендаций.
// Содержит ссылки на логгер и сервис рекомендаций.
type Handler struct {
	logger log.Factory
	svc    IRecommendationsService
}

// New - конструктор для создания нового обработчика Handler.
// Принимает логгер и сервис рекомендаций как параметры.
// Возвращает указатель на созданный обработчик.
func New(logger log.Factory, svc IRecommendationsService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}

}

// GetRecommendations - обработчик HTTP-запроса для получения всех рекомендаций.
// Запрашивает рекомендации из сервиса, обрабатывает возможные ошибки
// и возвращает результаты клиенту в формате JSON.
func (h *Handler) GetRecommendations(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	recommendations, err := h.svc.GetRecommendations()
	if err != nil {
		h.logger.Bg().Error("failed GetRecommendations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendations)
}

// GetRecommendationByID - обработчик HTTP-запроса для получения рекомендации по идентификатору.
// Запрашивает идентификатор из параметров запроса, получает рекомендацию из сервиса,
// обрабатывает возможные ошибки и возвращает результат клиенту в формате JSON.
func (h *Handler) GetRecommendationByID(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	id := c.Param("id")

	recommendation, err := h.svc.GetRecommendationByID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetRecommendationByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendation)
}

// GetRecommendationsByUserID - обработчик HTTP-запроса для получения рекомендаций по идентификатору пользователя.
// Запрашивает идентификатор пользователя из параметров запроса, получает рекомендации из сервиса,
// обрабатывает возможные ошибки и возвращает результаты клиенту в формате JSON.
func (h *Handler) GetRecommendationsByUserID(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	id := c.Param("id")

	recommendations, err := h.svc.GetRecommendationsByUserID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetRecommendationsByUserID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendations)
}
