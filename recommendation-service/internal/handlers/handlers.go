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
	//GetUsers() ([]model.User, error)
	//CreateUser(user model.User) (string, error)
	//UpdateUser(user model.User) (model.User, error)
	//GetUserByID(id string) (model.User, error)
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

//func (h *Handler) GetUsers(c *gin.Context) {
//	users, err := h.svc.GetUsers()
//	if err != nil {
//		h.logger.Bg().Error("failed GetUsers", zap.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, users)
//}
//
//func (h *Handler) CreateUser(c *gin.Context) {
//	var user model.User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		h.logger.Bg().Error("failed CreateUser", zap.Error(err))
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//
//	createdUserID, err := h.svc.CreateUser(user)
//	if err != nil {
//		h.logger.Bg().Error("failed CreateUser", zap.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusCreated, createdUserID)
//}
//
//func (h *Handler) GetUserByID(c *gin.Context) {
//	id := c.Param("id")
//
//	user, err := h.svc.GetUserByID(id)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//	}
//	c.JSON(http.StatusOK, user)
//}
