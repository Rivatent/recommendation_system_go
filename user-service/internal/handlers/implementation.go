package handlers

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"time"
	"user-service/internal/service"
	"user-service/pkg/log"
)

// Implementation включает в себя логгер для ведения журналов и указатель на HTTP-сервер, который будет обрабатывать запросы.
type Implementation struct {
	logger     log.Factory
	httpServer *http.Server
}

// NewServer инициализирует сервер, устанавливая логирование и параметр чтения, настраивает маршрутизацию.
func NewServer(logger log.Factory, svc *service.Service) *Implementation {
	return &Implementation{
		logger: logger,
		httpServer: &http.Server{
			ReadTimeout: 3 * time.Second,
			Handler:     newRouter(logger, svc),
		},
	}
}

// Run занимается фактическим запуском сервера, прослушиванием на заданном порту и логированием ошибок
func (s *Implementation) Run(_ context.Context) error {
	l, err := net.Listen("tcp", os.Getenv("HTTP_PORT"))
	if err != nil {
		s.logger.Bg().Error("in net.Listen", zap.Error(err))
		return err
	}

	if err = s.httpServer.Serve(l); err != nil {
		s.logger.Bg().Error("in httpServer.Serve", zap.Error(err))
		return err
	}

	return nil
}

// Stop - метод для остановки HTTP-сервера.
// Возвращает ошибку, если происходит ошибка при завершении работы сервера.
func (s *Implementation) Stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Bg().Error("in httpServer.Shutdown", zap.Error(err))
		return err
	}

	return nil
}
