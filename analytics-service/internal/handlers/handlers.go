package handlers

import (
	"analytics-service/internal/model"
	"analytics-service/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type IAnalyticsService interface {
	GetAnalytics() ([]model.Analytics, error)
}

type Handler struct {
	logger log.Factory
	svc    IAnalyticsService
}

func New(logger log.Factory, svc IAnalyticsService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}

}

func (h *Handler) GetAnalytics(c *gin.Context) {
	analytics, err := h.svc.GetAnalytics()
	if err != nil {
		h.logger.Bg().Error("failed GetAnalytics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analytics)
}
