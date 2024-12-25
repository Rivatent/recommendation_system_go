package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"user-service/internal/model"
	"user-service/internal/monitoring"
	"user-service/internal/validator"
	"user-service/log"
)

// IUserService интерфейс для взаимодействия с сервисом пользователей.
type IUserService interface {
	GetUsers() ([]model.User, error)
	CreateUser(user model.User) (string, error)
	UpdateUser(user model.User) (model.User, error)
	GetUserByID(id string) (model.User, error)
}

// Handler - структура, выполняющая функции обработки HTTP-запросов.
type Handler struct {
	logger log.Factory
	svc    IUserService
}

// New создает новый экземпляр Handler.
func New(logger log.Factory, svc IUserService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}
}

// UpdateUser обрабатывает HTTP-запрос на обновление пользователя.
func (h *Handler) UpdateUser(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Bg().Error("failed UpdateUser", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validator.Validate(user); err != nil {
		h.logger.Bg().Error("failed validation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.svc.UpdateUser(user)
	if err != nil {
		h.logger.Bg().Error("failed UpdateUser", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// GetUsers обрабатывает HTTP-запрос на получение списка пользователей.
func (h *Handler) GetUsers(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	users, err := h.svc.GetUsers()
	if err != nil {
		h.logger.Bg().Error("failed GetUsers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUser обрабатывает HTTP-запрос на создание нового пользователя.
func (h *Handler) CreateUser(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Bg().Error("failed CreateUser", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validator.Validate(user); err != nil {
		h.logger.Bg().Error("failed validation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUserID, err := h.svc.CreateUser(user)
	if err != nil {
		h.logger.Bg().Error("failed CreateUser", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdUserID})
}

// GetUserByID обрабатывает HTTP-запрос на получение пользователя по идентификатору.
func (h *Handler) GetUserByID(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	id := c.Param("id")

	user, err := h.svc.GetUserByID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetUserByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
