package service

import (
	"analytics-service/internal/model"
)

type IRepo interface {
	GetAnalyticsRepo() ([]model.Analytics, error)
	UpdateAnalyticsMsgRepo() error
}

type Service struct {
	repo IRepo
}

func New(repo IRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAnalytics() (analytics []model.Analytics, err error) {
	return s.repo.GetAnalyticsRepo()
}
