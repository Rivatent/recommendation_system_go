package service

import (
	"user-service/internal/model"
)

type IRepo interface {
	GetUsersRepo() ([]model.User, error)
	CreateUserRepo(user model.User) (string, error)
	UpdateUserRepo(user model.User) (model.User, error)
	GetUserByIDRepo(id string) (model.User, error)
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

func (s *Service) GetUsers() ([]model.User, error) {
	return s.repo.GetUsersRepo()
}

func (s *Service) CreateUser(user model.User) (string, error) {
	createdUserID, err := s.repo.CreateUserRepo(user)
	if err != nil {
		return createdUserID, err
	}
	user.ID = createdUserID
	updateMsg := map[string]interface{}{
		"user": user,
	}
	if err := s.KafkaProd.SendMessage(updateMsg, &s.KafkaProd.topicNew); err != nil {
		return createdUserID, err
	}

	return createdUserID, nil
}

func (s *Service) UpdateUser(user model.User) (model.User, error) {
	updatedUser, err := s.repo.UpdateUserRepo(user)
	if err != nil {
		return model.User{}, err
	}

	updateMessage := map[string]interface{}{
		"user": updatedUser,
	}
	if err := s.KafkaProd.SendMessage(updateMessage, &s.KafkaProd.topicUpdate); err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *Service) GetUserByID(id string) (model.User, error) {
	user, err := s.repo.GetUserByIDRepo(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
