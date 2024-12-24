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

type Implementation struct {
	logger log.Factory

	httpServer *http.Server
}

func NewServer(logger log.Factory, svc *service.Service) *Implementation {
	return &Implementation{
		logger: logger,
		httpServer: &http.Server{
			ReadTimeout: 3 * time.Second,
			Handler:     newRouter(logger, svc),
		},
	}
}

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

func (s *Implementation) Stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Bg().Error("in httpServer.Shutdown", zap.Error(err))
		return err
	}

	return nil
}
