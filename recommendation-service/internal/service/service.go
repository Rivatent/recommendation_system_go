package service

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"recommendation-service/internal/model"
	"time"
)

type IRepo interface {
	GetRecommendationsRepo() ([]model.Recommendation, error)
	GetRecommendationByIDRepo(id string) (model.Recommendation, error)
	GetRecommendationsByUserIDRepo(id string) ([]model.Recommendation, error)
	UserNewMsgRepo(newUser map[string]interface{}) error
	ProductNewMsgRepo(newProduct map[string]interface{}) error
	ProductUpdateMsgRepo(updatedProduct map[string]interface{}) error
}

type ICache interface {
	GetRecommendationByID(id string) (model.Recommendation, error)
	GetRecommendationsByUserID(id string) ([]model.Recommendation, error)
	SetRecommendationByID(id string, recommendation model.Recommendation, expiration time.Duration) error
	SetRecommendationsByUserID(id string, recommendations []model.Recommendation, expiration time.Duration) error
}

type Service struct {
	repo  IRepo
	cache ICache
}

func New(repo IRepo, cache ICache) *Service {

	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) GetRecommendations() (recommendations []model.Recommendation, err error) {
	return s.repo.GetRecommendationsRepo()
}

func (s *Service) GetRecommendationByID(id string) (model.Recommendation, error) {
	recommendation, err := s.cache.GetRecommendationByID(id)
	if err == nil && !errors.Is(err, redis.Nil) {
		return recommendation, nil
	}

	recommendation, err = s.repo.GetRecommendationByIDRepo(id)
	if err != nil {
		return model.Recommendation{}, err
	}

	err = s.cache.SetRecommendationByID(id, recommendation, 5*time.Minute)
	if err != nil {
		return recommendation, err
	}

	return recommendation, nil
}

func (s *Service) GetRecommendationsByUserID(id string) ([]model.Recommendation, error) {

	recommendations, err := s.cache.GetRecommendationsByUserID(id)
	if err == nil && !errors.Is(err, redis.Nil) {
		return recommendations, nil
	}

	recommendations, err = s.repo.GetRecommendationsByUserIDRepo(id)
	if err != nil {
		return nil, err
	}

	err = s.cache.SetRecommendationsByUserID(id, recommendations, 5*time.Minute)
	if err != nil {
		return recommendations, err
	}
	return recommendations, nil
}
