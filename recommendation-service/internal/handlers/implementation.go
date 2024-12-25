package handlers

import (
	"context"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"recommendation-service/internal/service"
	"recommendation-service/log"
	"time"
)

// Implementation - структура, представляющая сервер и его настройки.
// Содержит логгер для записи информации, а также HTTP-сервер.
type Implementation struct {
	logger log.Factory

	httpServer *http.Server
}

// NewServer - функция, создающая новый экземпляр Implementation.
// Принимает логгер и сервис в качестве параметров и возвращает указатель на
// структуру Implementation.
// logger - фабрика логирования для управления записями логов.
// svc - ссылка на сервис, который будет использоваться для обработки логики приложения.
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
// Принимает контекст для контроля выполнения.
// Возвращает ошибку, если происходит ошибка при прослушивании порта
// или при выполнении сервера.
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
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Bg().Error("in httpServer.Shutdown", zap.Error(err))
		return err
	}

	return nil
}
