package handlers

import (
	"context"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"product-service/internal/service"
	"product-service/pkg/log"
	"time"
)

// Implementation структура, представляющая собой реализацию HTTP-сервера
// для обработки запросов к продуктам. Содержит ссылку на экземпляр логгера
// и на HTTP-сервер.
type Implementation struct {
	logger log.Factory

	httpServer *http.Server
}

// NewServer создает новый экземпляр Implementation с заданным логгером и сервисом.
// Параметры:
//   - logger: экземпляр логгера, используемый для записи логов.
//   - svc: указатель на сервис, содержащий бизнес-логику для работы с продуктами.
//
// Возвращает указатель на новое значение Implementation.
func NewServer(logger log.Factory, svc *service.Service) *Implementation {
	return &Implementation{
		logger: logger,
		httpServer: &http.Server{
			ReadTimeout: 3 * time.Second,
			Handler:     newRouter(logger, svc),
		},
	}
}

// Run запускает HTTP-сервер и начинает слушать входящие соединения.
// Параметры:
//   - ctx: контекст для управления жизненным циклом сервера.
//
// Возвращает ошибку, если не удалось запустить сервер или получить соединение.
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

// Stop останавливает HTTP-сервер, завершая текущие соединения.
// Возвращает ошибку, если не удалось остановить сервер.
func (s *Implementation) Stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Bg().Error("in httpServer.Shutdown", zap.Error(err))
		return err
	}

	return nil
}
