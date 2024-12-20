package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"product-service/internal/model"
	"product-service/log"
)

type IProductService interface {
	GetProducts() ([]model.Product, error)
	CreateProduct(product model.Product) (string, error)
	UpdateProduct(product model.Product) (model.Product, error)
	GetProductByID(id string) (model.Product, error)
	DeleteProductByID(id string) error
}

type Handler struct {
	logger log.Factory
	svc    IProductService
}

func New(logger log.Factory, svc IProductService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}

}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Bg().Error("failed UpdateProduct", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	updatedProduct, err := h.svc.UpdateProduct(product)
	if err != nil {
		h.logger.Bg().Error("failed UpdateProduct", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (h *Handler) GetProducts(c *gin.Context) {
	products, err := h.svc.GetProducts()
	if err != nil {
		h.logger.Bg().Error("failed GetProducts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Bg().Error("failed CreateProduct", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdProductID, err := h.svc.CreateProduct(product)
	if err != nil {
		h.logger.Bg().Error("failed CreateProduct", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProductID)
}

func (h *Handler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.svc.GetProductByID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetProductByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) DeleteProductByID(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.DeleteProductByID(id)
	if err != nil {
		h.logger.Bg().Error("failed DeleteProductByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
