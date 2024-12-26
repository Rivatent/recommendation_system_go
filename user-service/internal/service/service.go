package service

import (
	"user-service/internal/model"
)

// IRepo - интерфейс для взаимодействия с репозиторием пользователей.
// Определяет методы для получения, создания, обновления и поиска пользователей по идентификатору.
type IRepo interface {
	GetUsersRepo() ([]model.User, error)
	CreateUserRepo(user model.User) (string, error)
	UpdateUserRepo(user model.User) (model.User, error)
	GetUserByIDRepo(id string) (model.User, error)
}

// Service - структура, представляющая бизнес-логику приложения слоя сервис.
// Содержит методы для работы с пользователями и ссылки на репозиторий и Kafka продюсер.
type Service struct {
	repo      IRepo
	KafkaProd IKafkaProducer
}

// New создает новый экземпляр Service.
// repo - реализация интерфейса IRepo для доступа к данным пользователей.
// kafkaProd - реализация интерфейса IKafkaProducer для отправки сообщений в Kafka.
// Возвращает указатель на новый экземпляр Service.
func New(repo IRepo, kafkaProd IKafkaProducer) *Service {
	return &Service{
		repo:      repo,
		KafkaProd: kafkaProd,
	}
}

// GetUsers возвращает список всех пользователей, получая их из репозитория.
// Возвращает срез пользователей и ошибку, если произошла ошибка.
func (s *Service) GetUsers() ([]model.User, error) {
	return s.repo.GetUsersRepo()
}

// CreateUser создает нового пользователя и отправляет сообщение о создании в Kafka.
// user - объект пользователя, которого нужно создать.
// Возвращает идентификатор созданного пользователя и ошибку, если произошла ошибка.
func (s *Service) CreateUser(user model.User) (string, error) {
	createdUserID, err := s.repo.CreateUserRepo(user)
	if err != nil {
		return createdUserID, err
	}
	user.ID = createdUserID
	updateMsg := map[string]interface{}{
		"user": user,
	}
	if err := s.KafkaProd.SendMessage(updateMsg, s.KafkaProd.TopicNew()); err != nil {
		return createdUserID, err
	}

	return createdUserID, nil
}

// UpdateUser обновляет информацию о пользователе и отправляет сообщение об обновлении в Kafka.
// user - объект с обновленной информацией о пользователе.
// Возвращает обновленный объект пользователя и ошибку, если произошла ошибка.
func (s *Service) UpdateUser(user model.User) (model.User, error) {
	updatedUser, err := s.repo.UpdateUserRepo(user)
	if err != nil {
		return model.User{}, err
	}

	updateMessage := map[string]interface{}{
		"user": updatedUser,
	}
	if err := s.KafkaProd.SendMessage(updateMessage, s.KafkaProd.TopicUpdate()); err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

// GetUserByID ищет пользователя по его идентификатору.
// id - строка, представляющая идентификатор пользователя.
// Возвращает объект пользователя и ошибку, если пользователь не найден или произошла ошибка.
func (s *Service) GetUserByID(id string) (model.User, error) {
	user, err := s.repo.GetUserByIDRepo(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
