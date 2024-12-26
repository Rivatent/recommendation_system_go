package service

import (
	"analytics-service/internal/model"
)

// IRepo - интерфейс, определяющий методы для доступа к данным аналитики.
// Этот интерфейс позволяет получать данные аналитики и обновлять их в репозитории.
type IRepo interface {
	GetAnalyticsRepo() ([]model.Analytics, error)
	UpdateAnalyticsMsgRepo() error
}

// Service - структура, представляющая сервис для работы с аналитикой.
// Она содержит ссылку на интерфейс IRepo для доступа и управления данными аналитики.
type Service struct {
	repo IRepo
}

// New - функция для создания нового экземпляра Service.
// Она принимает интерфейс IRepo как аргумент и возвращает указатель на новый экземпляр Service.
//
// Параметры:
// - repo IRepo: интерфейс репозитория для взаимодействия с данными аналитики.
//
// Возвращаемое значение:
// - *Service: указатель на структуру Service, готовую к использованию.
func New(repo IRepo) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAnalytics - метод для получения аналитических данных.
// Он использует метод GetAnalyticsRepo интерфейса IRepo для извлечения данных и
// возвращает срез аналитики и ошибку, если такая имеется.
//
// Возвращаемые значения:
// - []model.Analytics: срез структур Analytics, представляющий данные аналитики.
// - error: ошибка, если возникла проблема при получении данных.
func (s *Service) GetAnalytics() (analytics []model.Analytics, err error) {
	return s.repo.GetAnalyticsRepo()
}
