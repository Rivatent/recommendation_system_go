package service

import "product-service/internal/repository"

type IRepo interface {
	//GetUsersRepo() ([]repository.User, error)
	//CreateUserRepo(user repository.User) (repository.User, error)
	//UpdateUserRepo(user repository.User) (repository.User, error)
	//GetUserByIDRepo(id int) (repository.User, error)
	GetProductsRepo() ([]repository.Product, error)
}

type Service struct {
	repo IRepo
}

func New(repo IRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetProducts() ([]repository.Product, error) {
	return s.repo.GetProductsRepo()
}

//func (s *Service) GetUsers() ([]repository.User, error) {
//	return s.repo.GetUsersRepo()
//}
//
//func (s *Service) CreateUser(user repository.User) (repository.User, error) {
//	createdUser, err := s.repo.CreateUserRepo(user)
//	if err != nil {
//		return repository.User{}, err
//	}
//
//	return createdUser, nil
//}
//
//func (s *Service) UpdateUser(user repository.User) (repository.User, error) {
//	updatedUser, err := s.repo.UpdateUserRepo(user)
//	if err != nil {
//		return repository.User{}, err
//	}
//
//	return updatedUser, nil
//}
//
//func (s *Service) GetUserByID(id int) (repository.User, error) {
//	user, err := s.repo.GetUserByIDRepo(id)
//	if err != nil {
//		return repository.User{}, err
//	}
//	return user, nil
//}
