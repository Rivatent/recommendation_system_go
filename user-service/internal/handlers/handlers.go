package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-service/internal/repository"
)

type IUserService interface {
	GetUsers() ([]repository.User, error)
	CreateUser(user repository.User) (repository.User, error)
	UpdateUser(user repository.User) (repository.User, error)
	GetUserByID(id int) (repository.User, error)
}

type Handler struct {
	svc IUserService
}

func New(svc IUserService) *Handler {
	return &Handler{
		svc: svc,
	}

}

func (h *Handler) UpdateUser(c *gin.Context) {
	var user repository.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	updatedUser, err := h.svc.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.svc.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user repository.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdUser, err := h.svc.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	}
	user, err := h.svc.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}
