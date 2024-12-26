package service

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"recommendation-service/internal/model"
	"time"
)

// IRepo является интерфейсом, который определяет методы для взаимодействия с репозиторием рекомендаций.
type IRepo interface {
	GetRecommendationsRepo() ([]model.Recommendation, error)
	GetRecommendationByIDRepo(id string) (model.Recommendation, error)
	GetRecommendationsByUserIDRepo(id string) ([]model.Recommendation, error)
	UserNewMsgRepo(newUser map[string]interface{}) error
	ProductNewMsgRepo(newProduct map[string]interface{}) error
	ProductUpdateMsgRepo(updatedProduct map[string]interface{}) error
}

// ICache является интерфейсом, который определяет методы для взаимодействия с кэшем.
// Он предоставляет методы для получения и установки рекомендаций в кэш, а также для
// установки времени истечения.
type ICache interface {
	GetRecommendationByID(id string) (model.Recommendation, error)
	GetRecommendationsByUserID(id string) ([]model.Recommendation, error)
	SetRecommendationByID(id string, recommendation model.Recommendation, expiration time.Duration) error
	SetRecommendationsByUserID(id string, recommendations []model.Recommendation, expiration time.Duration) error
}

// Service представляет собой логику для управления рекомендациями.
// Он использует репозиторий (IRepo) для доступа к данным и кэш (ICache)
// для повышения производительности операций чтения.
type Service struct {
	repo  IRepo
	cache ICache
}

// New создает новый экземпляр Service с заданными репозиторием и кэшем.
// Возвращает указатель на Service.
func New(repo IRepo, cache ICache) *Service {

	return &Service{
		repo:  repo,
		cache: cache,
	}
}

// GetRecommendations возвращает все рекомендации, получая их из репозитория.
// В случае ошибки возвращает ошибку.
func (s *Service) GetRecommendations() (recommendations []model.Recommendation, err error) {
	return s.repo.GetRecommendationsRepo()
}

// GetRecommendationByID возвращает рекомендацию по ее идентификатору.
// Сначала проверяет наличие в кэше, если не найдено - загружает из репозитория.
// При успешной загрузке сохраняет рекомендацию в кэше на 5 минут.
// Возвращает рекомендацию и возможную ошибку.
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

// GetRecommendationsByUserID возвращает рекомендации для пользователя по его идентификатору.
// Сначала проверяет наличие в кэше, если не найдено - загружает из репозитория.
// При успешной загрузке сохраняет рекомендации в кэше на 5 минут.
// Возвращает список рекомендаций и возможную ошибку.
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
