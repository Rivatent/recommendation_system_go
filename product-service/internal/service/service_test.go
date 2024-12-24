package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"product-service/internal/model"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetProductsRepo() ([]model.Product, error) {
	args := m.Called()
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockRepo) CreateProductRepo(Product model.Product) (string, error) {
	args := m.Called(Product)
	return args.String(0), args.Error(1)
}

func (m *MockRepo) UpdateProductRepo(Product model.Product) (model.Product, error) {
	args := m.Called(Product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockRepo) GetProductByIDRepo(id string) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockRepo) DeleteProductByIDRepo(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) SendMessage(message interface{}, topic *string) error {
	args := m.Called(message, topic)
	return args.Error(0)
}

func (m *MockKafkaProducer) TopicNew() *string {
	args := m.Called()
	return args.Get(0).(*string)
}

func (m *MockKafkaProducer) TopicUpdate() *string {
	args := m.Called()
	return args.Get(0).(*string)
}

func TestGetProducts(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	products := []model.Product{{ID: "1", Name: "Product 1"}, {ID: "2", Name: "Product 2"}}
	mockRepo.On("GetProductsRepo").Return(products, nil)

	result, err := service.GetProducts()

	assert.NoError(t, err)
	assert.Equal(t, products, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafkaProd := new(MockKafkaProducer)
	service := New(mockRepo, mockKafkaProd)

	product := model.Product{Name: "New Product"}
	expectedID := "123"
	expectedTopic := "TopicNew"

	mockRepo.On("CreateProductRepo", product).Return(expectedID, nil)
	mockKafkaProd.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
	mockKafkaProd.On("TopicNew").Return(&expectedTopic)

	id, err := service.CreateProduct(product)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)
	mockKafkaProd.AssertExpectations(t)
}

func TestCreateProduct_Error(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafkaProd := new(MockKafkaProducer)
	service := New(mockRepo, mockKafkaProd)

	product := model.Product{Name: "New Product"}
	mockRepo.On("CreateProductRepo", product).Return("", assert.AnError)

	id, err := service.CreateProduct(product)

	assert.Error(t, err)
	assert.Equal(t, "", id)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafkaProd := new(MockKafkaProducer)
	service := New(mockRepo, mockKafkaProd)
	expectedUpdateTopic := "TopicUpdate"

	product := model.Product{ID: "123", Name: "Updated Product"}
	mockRepo.On("UpdateProductRepo", product).Return(product, nil)
	mockKafkaProd.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
	mockKafkaProd.On("TopicUpdate").Return(&expectedUpdateTopic)

	updatedProduct, err := service.UpdateProduct(product)

	assert.NoError(t, err)
	assert.Equal(t, product, updatedProduct)
	mockRepo.AssertExpectations(t)
	mockKafkaProd.AssertExpectations(t)
}

func TestUpdateProduct_Error(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafkaProd := new(MockKafkaProducer)
	service := New(mockRepo, mockKafkaProd)

	product := model.Product{ID: "123", Name: "Updated Product"}
	mockRepo.On("UpdateProductRepo", product).Return(model.Product{}, assert.AnError)

	updatedProduct, err := service.UpdateProduct(product)

	assert.Error(t, err)
	assert.Equal(t, model.Product{}, updatedProduct)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	product := model.Product{ID: "123", Name: "Product by ID"}
	mockRepo.On("GetProductByIDRepo", "123").Return(product, nil)

	result, err := service.GetProductByID("123")

	assert.NoError(t, err)
	assert.Equal(t, product, result)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_Error(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	mockRepo.On("GetProductByIDRepo", "123").Return(model.Product{}, assert.AnError)

	result, err := service.GetProductByID("123")

	assert.Error(t, err)
	assert.Equal(t, model.Product{}, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProductByID_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	mockRepo.On("DeleteProductByIDRepo", "123").Return(nil)

	err := service.DeleteProductByID("123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProductByID_Error(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	mockRepo.On("DeleteProductByIDRepo", "123").Return(assert.AnError)

	err := service.DeleteProductByID("123")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
