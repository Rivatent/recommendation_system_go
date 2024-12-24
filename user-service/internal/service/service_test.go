package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"user-service/internal/model"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetUsersRepo() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockRepo) CreateUserRepo(user model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockRepo) UpdateUserRepo(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockRepo) GetUserByIDRepo(id string) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
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

func TestGetUsers_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	expectedUsers := []model.User{
		{ID: "1", Username: "user1", Email: "user1@example.com"},
		{ID: "2", Username: "user2", Email: "user2@example.com"},
	}

	mockRepo.On("GetUsersRepo").Return(expectedUsers, nil)

	users, err := service.GetUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_Error(t *testing.T) {
	mockRepo := new(MockRepo)
	service := New(mockRepo, nil)

	mockRepo.On("GetUsersRepo").Return([]model.User{}, errors.New("database error"))

	_, err := service.GetUsers()

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafka := new(MockKafkaProducer)
	service := New(mockRepo, mockKafka)

	newUser := model.User{Username: "user1", Email: "user1@example.com", Password: "password"}
	createdUserID := "123"
	topicNew := "topic_new"

	mockRepo.On("CreateUserRepo", newUser).Return(createdUserID, nil)
	mockKafka.On("TopicNew").Return(&topicNew)
	mockKafka.On("SendMessage", mock.Anything, &topicNew).Return(nil)

	id, err := service.CreateUser(newUser)

	assert.NoError(t, err)
	assert.Equal(t, createdUserID, id)

	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestCreateUser_RepoError(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafka := new(MockKafkaProducer)
	service := New(mockRepo, mockKafka)

	newUser := model.User{Username: "user1", Email: "user1@example.com", Password: "password"}
	expectedError := errors.New("repository error")

	mockRepo.On("CreateUserRepo", newUser).Return("", expectedError)

	id, err := service.CreateUser(newUser)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, id)

	mockRepo.AssertExpectations(t)
	mockKafka.AssertNotCalled(t, "SendMessage")
}

func TestCreateUser_KafkaError(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafka := new(MockKafkaProducer)
	service := New(mockRepo, mockKafka)

	newUser := model.User{Username: "user1", Email: "user1@example.com", Password: "password"}
	createdUserID := "123"
	topicNew := "topic_new"
	expectedError := errors.New("kafka error")

	mockRepo.On("CreateUserRepo", newUser).Return(createdUserID, nil)
	mockKafka.On("TopicNew").Return(&topicNew)
	mockKafka.On("SendMessage", mock.Anything, &topicNew).Return(expectedError)

	id, err := service.CreateUser(newUser)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, createdUserID, id)

	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafka := new(MockKafkaProducer)
	service := New(mockRepo, mockKafka)

	updatedUser := model.User{ID: "123", Username: "user1", Email: "user1@example.com"}
	topicUpdate := "topic_update"

	mockRepo.On("UpdateUserRepo", updatedUser).Return(updatedUser, nil)
	mockKafka.On("TopicUpdate").Return(&topicUpdate)
	mockKafka.On("SendMessage", mock.Anything, &topicUpdate).Return(nil)

	result, err := service.UpdateUser(updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, result)

	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafka := new(MockKafkaProducer)
	service := New(mockRepo, mockKafka)

	expectedUser := model.User{ID: "123", Username: "user1", Email: "user1@example.com"}

	mockRepo.On("GetUserByIDRepo", "123").Return(expectedUser, nil)

	user, err := service.GetUserByID("123")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	mockKafka := new(MockKafkaProducer)
	service := New(mockRepo, mockKafka)

	mockRepo.On("GetUserByIDRepo", "123").Return(model.User{}, errors.New("user not found"))

	user, err := service.GetUserByID("123")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, model.User{}, user)

	mockRepo.AssertExpectations(t)
}
