package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"recommendation-service/internal/model"
	"recommendation-service/log"
)

type IRecommendationsService interface {
	GetRecommendations() ([]model.Recommendation, error)
	GetRecommendationByID(id string) (model.Recommendation, error)
	GetRecommendationsByUserID(id string) ([]model.Recommendation, error)
}

type Handler struct {
	logger log.Factory
	svc    IRecommendationsService
}

func New(logger log.Factory, svc IRecommendationsService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}

}

func (h *Handler) GetRecommendations(c *gin.Context) {
	recommendations, err := h.svc.GetRecommendations()
	if err != nil {
		h.logger.Bg().Error("failed GetRecommendations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendations)
}

func (h *Handler) GetRecommendationByID(c *gin.Context) {
	id := c.Param("id")

	recommendation, err := h.svc.GetRecommendationByID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetRecommendationByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendation)
}

func (h *Handler) GetRecommendationsByUserID(c *gin.Context) {
	id := c.Param("id")

	recommendations, err := h.svc.GetRecommendationsByUserID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetRecommendationsByUserID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendations)
}
