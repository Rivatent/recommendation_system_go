package handlers

import (
	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetUsers()
}

type Handler struct {
	svc IUserService
}

func New(svc IUserService) *Handler {
	return &Handler{
		svc: svc,
	}

}

func (h *Handler) GetUsers(c *gin.Context) {
	h.svc.GetUsers()
}

func (h *Handler) CreateUser(c *gin.Context) {
}
