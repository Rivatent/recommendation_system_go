package service

import (
	"product-service/internal/model"
)

// IRepo - интерфейс для управления репозиторием продуктов.
// Определяет методы для работы с продуктами,
// включая получение, создание, обновление и удаление.
type IRepo interface {
	GetProductsRepo() ([]model.Product, error)
	CreateProductRepo(Product model.Product) (string, error)
	UpdateProductRepo(Product model.Product) (model.Product, error)
	GetProductByIDRepo(id string) (model.Product, error)
	DeleteProductByIDRepo(id string) error
}

// Service - структура, представляющая сервис управления продуктами.
// Содержит репозиторий и Kafka продюсер для выполнения операций с продуктами.
type Service struct {
	repo      IRepo
	KafkaProd IKafkaProducer
}

// New создает новый экземпляр сервиса.
// Принимает в качестве параметров интерфейс репозитория и Kafka продюсера.
func New(repo IRepo, kafkaProd IKafkaProducer) *Service {
	return &Service{
		repo:      repo,
		KafkaProd: kafkaProd,
	}
}

// GetProducts получает все продукты из репозитория.
func (s *Service) GetProducts() ([]model.Product, error) {
	return s.repo.GetProductsRepo()
}

// CreateProduct создает новый продукт и отправляет сообщение в Kafka о создании продукта.
func (s *Service) CreateProduct(Product model.Product) (string, error) {
	createdProductID, err := s.repo.CreateProductRepo(Product)
	if err != nil {
		return createdProductID, err
	}
	Product.ID = createdProductID
	updateMsg := map[string]interface{}{
		"product": Product,
	}
	if err := s.KafkaProd.SendMessage(updateMsg, s.KafkaProd.TopicNew()); err != nil {
		return createdProductID, err
	}

	return createdProductID, nil
}

// UpdateProduct обновляет продукт и отправляет сообщение в Kafka о его обновлении.
func (s *Service) UpdateProduct(Product model.Product) (model.Product, error) {
	updatedProduct, err := s.repo.UpdateProductRepo(Product)
	if err != nil {
		return model.Product{}, err
	}
	updateMessage := map[string]interface{}{
		"product": updatedProduct,
	}
	if err := s.KafkaProd.SendMessage(updateMessage, s.KafkaProd.TopicUpdate()); err != nil {
		return updatedProduct, err
	}

	return updatedProduct, nil
}

// GetProductByID получает продукт по его уникальному идентификатору.
func (s *Service) GetProductByID(id string) (model.Product, error) {
	Product, err := s.repo.GetProductByIDRepo(id)
	if err != nil {
		return model.Product{}, err
	}

	return Product, nil
}

// DeleteProductByID удаляет продукт по его уникальному идентификатору.
func (s *Service) DeleteProductByID(id string) error {
	err := s.repo.DeleteProductByIDRepo(id)
	if err != nil {
		return err
	}

	return nil
}
