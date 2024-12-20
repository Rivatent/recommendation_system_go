package service

import "recommendation-service/internal/model"

type IRepo interface {
	GetRecommendationsRepo() ([]model.Recommendation, error)
	//GetUsersRepo() ([]model.User, error)
	//CreateUserRepo(user model.User) (string, error)
	//UpdateUserRepo(user model.User) (model.User, error)
	//GetUserByIDRepo(id string) (model.User, error)
}

type Service struct {
	repo IRepo
}

func New(repo IRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetRecommendations() (recommendations []model.Recommendation, err error) {
	return s.repo.GetRecommendationsRepo()
}

//func (s *Service) GetUsers() ([]model.User, error) {
//	return s.repo.GetUsersRepo()
//}
//
//func (s *Service) CreateUser(user model.User) (string, error) {
//	createdUserID, err := s.repo.CreateUserRepo(user)
//	if err != nil {
//		return createdUserID, err
//	}
//
//	updateMsg := map[string]interface{}{
//		"event": "user_created",
//		"user":  user,
//		"id":    createdUserID,
//	}
//	if err := s.KafkaProd.SendMessage(updateMsg); err != nil {
//		return createdUserID, err
//	}
//
//	return createdUserID, nil
//}
//
//func (s *Service) UpdateUser(user model.User) (model.User, error) {
//	updatedUser, err := s.repo.UpdateUserRepo(user)
//	if err != nil {
//		return model.User{}, err
//	}
//
//	updateMessage := map[string]interface{}{
//		"event": "user_updated",
//		"user":  updatedUser,
//	}
//	if err := s.KafkaProd.SendMessage(updateMessage); err != nil {
//		return updatedUser, err
//	}
//
//	return updatedUser, nil
//}
//
//func (s *Service) GetUserByID(id string) (model.User, error) {
//	user, err := s.repo.GetUserByIDRepo(id)
//	if err != nil {
//		return model.User{}, err
//	}
//	return user, nil
//}
