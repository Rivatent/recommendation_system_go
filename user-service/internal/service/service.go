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
	repo IRepo
}

func New(repo IRepo) *Service {
	return &Service{
		repo: repo,
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

	return createdUserID, nil
}

func (s *Service) UpdateUser(user model.User) (model.User, error) {
	updatedUser, err := s.repo.UpdateUserRepo(user)
	if err != nil {
		return model.User{}, err
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
