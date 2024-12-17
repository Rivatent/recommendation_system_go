package service

type IRepo interface {
	GetUsersRepo()
}

type Service struct {
	repo IRepo
}

func New(repo IRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUsers() {
	s.repo.GetUsersRepo()
}
