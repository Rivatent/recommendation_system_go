package service

import (
	"product-service/internal/model"
)

type IRepo interface {
	GetProductsRepo() ([]model.Product, error)
	CreateProductRepo(Product model.Product) (string, error)
	UpdateProductRepo(Product model.Product) (model.Product, error)
	GetProductByIDRepo(id string) (model.Product, error)
	DeleteProductByIDRepo(id string) error
}

type Service struct {
	repo      IRepo
	KafkaProd *KafkaProducer
}

func New(repo IRepo, kafkaProd *KafkaProducer) *Service {
	return &Service{
		repo:      repo,
		KafkaProd: kafkaProd,
	}
}

func (s *Service) GetProducts() ([]model.Product, error) {
	return s.repo.GetProductsRepo()
}

func (s *Service) CreateProduct(Product model.Product) (string, error) {
	createdProductID, err := s.repo.CreateProductRepo(Product)
	if err != nil {
		return createdProductID, err
	}

	updateMsg := map[string]interface{}{
		"event":   "Product_created",
		"Product": Product,
		"id":      createdProductID,
	}
	if err := s.KafkaProd.SendMessage(updateMsg); err != nil {
		return createdProductID, err
	}

	return createdProductID, nil
}

func (s *Service) UpdateProduct(Product model.Product) (model.Product, error) {
	updatedProduct, err := s.repo.UpdateProductRepo(Product)
	if err != nil {
		return model.Product{}, err
	}

	updateMessage := map[string]interface{}{
		"event":   "Product_updated",
		"Product": updatedProduct,
	}
	if err := s.KafkaProd.SendMessage(updateMessage); err != nil {
		return updatedProduct, err
	}

	return updatedProduct, nil
}

func (s *Service) GetProductByID(id string) (model.Product, error) {
	Product, err := s.repo.GetProductByIDRepo(id)
	if err != nil {
		return model.Product{}, err
	}
	return Product, nil
}

func (s *Service) DeleteProductByID(id string) error {
	err := s.repo.DeleteProductByIDRepo(id)
	if err != nil {
		return err
	}
	return nil
}
