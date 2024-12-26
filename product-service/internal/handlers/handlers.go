package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"product-service/internal/model"
	"product-service/internal/monitoring"
	"product-service/internal/validator"
	"product-service/log"
	"time"
)

// IProductService интерфейс, предоставляющий методы для работы с продуктами.
// Реализация этого интерфейса обрабатывает запросы, связанные с продуктами,
// такие как получение, создание, обновление и удаление продуктов.
type IProductService interface {
	GetProducts() ([]model.Product, error)
	CreateProduct(product model.Product) (string, error)
	UpdateProduct(product model.Product) (model.Product, error)
	GetProductByID(id string) (model.Product, error)
	DeleteProductByID(id string) error
}

// Handler структура, содержащая логгер и сервис продуктов.
// Обработчики в этом пакете используют данный Handler для выполнения операций
// с продуктами
type Handler struct {
	logger log.Factory
	svc    IProductService
}

// New создает новый экземпляр Handler с заданным логгером и сервисом продуктов.
// Параметры:
//   - logger: экземпляр логгера для записи сообщений.
//   - svc: экземпляр службы продуктов для выполнения бизнес-логики.
func New(logger log.Factory, svc IProductService) *Handler {
	return &Handler{
		logger: logger,
		svc:    svc,
	}

}

// UpdateProduct обрабатывает HTTP запрос на обновление продукта.
// Обновляет продукт и возвращает его обновленную информацию в формате JSON.
// Параметры:
//   - c: контекст Gin, содержащий информацию о запросе и ответе.
func (h *Handler) UpdateProduct(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Bg().Error("failed UpdateProduct", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := validator.Validate(product); err != nil {
		h.logger.Bg().Error("failed validation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// GetProducts обрабатывает HTTP запрос на получение всех продуктов.
// Возвращает список продуктов в формате JSON.
// Параметры:
//   - c: контекст Gin, содержащий информацию о запросе и ответе.
func (h *Handler) GetProducts(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	products, err := h.svc.GetProducts()
	if err != nil {
		h.logger.Bg().Error("failed GetProducts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateProduct обрабатывает HTTP запрос на создание нового продукта.
// Создает продукт и возвращает его уникальный идентификатор в формате JSON.
// Параметры:
//   - c: контекст Gin, содержащий информацию о запросе и ответе.
func (h *Handler) CreateProduct(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Bg().Error("failed CreateProduct", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validator.Validate(product); err != nil {
		h.logger.Bg().Error("failed validation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProductID, err := h.svc.CreateProduct(product)
	if err != nil {
		h.logger.Bg().Error("failed CreateProduct", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdProductID})
}

// GetProductByID обрабатывает HTTP запрос на получение продукта по его уникальному идентификатору.
// Возвращает продукт в формате JSON.
// Параметры:
//   - c: контекст Gin, содержащий информацию о запросе и ответе.
func (h *Handler) GetProductByID(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	id := c.Param("id")

	product, err := h.svc.GetProductByID(id)
	if err != nil {
		h.logger.Bg().Error("failed GetProductByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProductByID обрабатывает HTTP запрос на удаление продукта по его уникальному идентификатору.
// Удаляет продукт и возвращает статус 204 No Content.
// Параметры:
//   - c: контекст Gin, содержащий информацию о запросе и ответе.
func (h *Handler) DeleteProductByID(c *gin.Context) {
	start := time.Now()
	defer monitoring.CollectMetrics(start, c)

	id := c.Param("id")

	err := h.svc.DeleteProductByID(id)
	if err != nil {
		h.logger.Bg().Error("failed DeleteProductByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
