package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"recommendation-service/internal/model"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetRecommendationsRepo() ([]model.Recommendation, error) {
	args := m.Called()
	return args.Get(0).([]model.Recommendation), args.Error(1)
}

func (m *MockRepo) GetRecommendationByIDRepo(id string) (model.Recommendation, error) {
	args := m.Called(id)
	return args.Get(0).(model.Recommendation), args.Error(1)
}

func (m *MockRepo) GetRecommendationsByUserIDRepo(id string) ([]model.Recommendation, error) {
	args := m.Called(id)
	return args.Get(0).([]model.Recommendation), args.Error(1)
}

func (m *MockRepo) UserNewMsgRepo(newUser map[string]interface{}) error {
	args := m.Called(newUser)
	return args.Error(0)
}

func (m *MockRepo) ProductNewMsgRepo(newProduct map[string]interface{}) error {
	args := m.Called(newProduct)
	return args.Error(0)
}

func (m *MockRepo) ProductUpdateMsgRepo(updatedProduct map[string]interface{}) error {
	args := m.Called(updatedProduct)
	return args.Error(0)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetRecommendationByID(id string) (model.Recommendation, error) {
	args := m.Called(id)
	return args.Get(0).(model.Recommendation), args.Error(1)
}

func (m *MockCache) GetRecommendationsByUserID(id string) ([]model.Recommendation, error) {
	args := m.Called(id)
	return args.Get(0).([]model.Recommendation), args.Error(1)
}

func (m *MockCache) SetRecommendationByID(id string, recommendation model.Recommendation, expiration time.Duration) error {
	args := m.Called(id, recommendation, expiration)
	return args.Error(0)
}

func (m *MockCache) SetRecommendationsByUserID(id string, recommendations []model.Recommendation, expiration time.Duration) error {
	args := m.Called(id, recommendations, expiration)
	return args.Error(0)
}

func TestService_GetRecommendations(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)
	service := New(mockRepo, mockCache)

	expectedRecommendations := []model.Recommendation{
		{ID: "1", ProductID: "product1", UserID: "user1", Score: 4.5},
		{ID: "2", ProductID: "product2", UserID: "user2", Score: 4.7},
	}

	mockRepo.On("GetRecommendationsRepo").Return(expectedRecommendations, nil)

	recommendations, err := service.GetRecommendations()

	assert.NoError(t, err)
	assert.Equal(t, expectedRecommendations, recommendations)
	mockRepo.AssertExpectations(t)
}

func TestService_GetRecommendationByID(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)
	service := New(mockRepo, mockCache)

	expectedRecommendation := model.Recommendation{ID: "1", ProductID: "product1", UserID: "user1", Score: 4.5}

	mockCache.On("GetRecommendationByID", "1").Return(expectedRecommendation, nil)

	recommendation, err := service.GetRecommendationByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expectedRecommendation, recommendation)
	mockCache.AssertExpectations(t)

	mockCache.On("GetRecommendationByID", "2").Return(model.Recommendation{}, redis.Nil)
	mockRepo.On("GetRecommendationByIDRepo", "2").Return(expectedRecommendation, nil)
	mockCache.On("SetRecommendationByID", "2", expectedRecommendation, 5*time.Minute).Return(nil)

	recommendation, err = service.GetRecommendationByID("2")

	assert.NoError(t, err)
	assert.Equal(t, expectedRecommendation, recommendation)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestService_GetRecommendationsByUserID(t *testing.T) {
	mockRepo := new(MockRepo)
	mockCache := new(MockCache)
	service := New(mockRepo, mockCache)

	expectedRecommendations := []model.Recommendation{
		{ID: "1", ProductID: "product1", UserID: "user1", Score: 4.5},
		{ID: "2", ProductID: "product2", UserID: "user1", Score: 4.7},
	}

	mockCache.On("GetRecommendationsByUserID", "user1").Return(expectedRecommendations, nil)

	recommendations, err := service.GetRecommendationsByUserID("user1")

	assert.NoError(t, err)
	assert.Equal(t, expectedRecommendations, recommendations)
	mockCache.AssertExpectations(t)

	mockCache.On("GetRecommendationsByUserID", "user2").Return([]model.Recommendation{}, redis.Nil)
	mockRepo.On("GetRecommendationsByUserIDRepo", "user2").Return(expectedRecommendations, nil)
	mockCache.On("SetRecommendationsByUserID", "user2", expectedRecommendations, 5*time.Minute).Return(nil)

	recommendations, err = service.GetRecommendationsByUserID("user2")

	assert.NoError(t, err)
	assert.Equal(t, expectedRecommendations, recommendations)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
