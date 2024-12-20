package service

import "recommendation-service/internal/model"

type IRepo interface {
	GetRecommendationsRepo() ([]model.Recommendation, error)
	GetRecommendationByIDRepo(id string) (model.Recommendation, error)
	GetRecommendationsByUserIDRepo(id string) ([]model.Recommendation, error)
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

func (s *Service) GetRecommendationByID(id string) (model.Recommendation, error) {
	recommendation, err := s.repo.GetRecommendationByIDRepo(id)
	if err != nil {
		return model.Recommendation{}, err
	}
	return recommendation, nil
}

func (s *Service) GetRecommendationsByUserID(id string) ([]model.Recommendation, error) {
	recommendations, err := s.repo.GetRecommendationsByUserIDRepo(id)
	if err != nil {
		return nil, err
	}
	return recommendations, nil
}
