package handlers

import (
	"analytics-service/internal/service"
	"analytics-service/pkg/log"
	"context"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"time"
)

// Implementation - структура, представляющая реализацию HTTP-сервера.
// Содержит логгер и экземпляр http.Server для обработки входящих HTTP-запросов.
type Implementation struct {
	logger log.Factory

	httpServer *http.Server
}

// NewServer - функция для создания нового экземпляра Implementation.
// Принимает логгер и сервис аналитики в качестве аргументов.
// Возвращает указатель на созданный экземпляр Implementation.
func NewServer(logger log.Factory, svc *service.Service) *Implementation {
	return &Implementation{
		logger: logger,
		httpServer: &http.Server{
			ReadTimeout: 3 * time.Second,
			Handler:     newRouter(logger, svc),
		},
	}
}

// Run - метод для запуска HTTP-сервера.
// Принимает контекст для управления жизненным циклом.
// Возвращает ошибку, если возникли проблемы при запуске сервера.
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

// Stop - метод для корректного завершения работы HTTP-сервера.
// Освобождает все ресурсы, связанные с сервером.
// Возвращает ошибку, если возникли проблемы во время завершения.
func (s *Implementation) Stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Bg().Error("in httpServer.Shutdown", zap.Error(err))
		return err
	}

	return nil
}
