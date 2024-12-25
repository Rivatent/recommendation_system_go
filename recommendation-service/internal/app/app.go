package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"recommendation-service/cache"
	"recommendation-service/internal/closer"
	"recommendation-service/internal/handlers"
	"recommendation-service/internal/repository"
	"recommendation-service/internal/service"
	"recommendation-service/log"
	"sync"
)

// IRunner представляет интерфейс для запуска и остановки компонентов приложения.
type IRunner interface {
	Run(ctx context.Context) error
	Stop() error
}

// App приложение, которое включает сервер HTTP и kafka-консюмер.
type App struct {
	serverHTTP    IRunner
	kafkaConsumer IRunner
}

// New создает новое приложение с инициализацией необходимых компонентов, включая:
// логирование, базу данных, HTTP-сервер и Kafka-консюмер.
func New() (*App, error) {
	l := log.InitLogger().With(zap.String("app", "recommendation-service"))

	appLogger := log.NewFactory(l)

	db := repository.New()
	closer.Add(db.Close)

	redisCache := cache.NewRedisCache()
	closer.Add(redisCache.Close)

	svc := service.New(db, redisCache)

	httpSrv := handlers.NewServer(appLogger, svc)
	kafkaConsumer := service.NewKafkaConsumer(appLogger, db)

	return &App{
		serverHTTP:    httpSrv,
		kafkaConsumer: kafkaConsumer,
	}, nil
}

// Run запускает все компоненты приложения в отдельных горутинах.
// Ожидает сигнал завершения и после этого производит остановку всех компонентов.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	run := func(runner IRunner) {
		wg.Add(1)
		defer wg.Done()

		err := runner.Run(ctx)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(err)
		}
	}

	go run(a.serverHTTP)
	closer.Add(a.serverHTTP.Stop)

	go run(a.kafkaConsumer)
	closer.Add(a.kafkaConsumer.Stop)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan

	wg.Wait()
	return nil
}

// Stop останавливает выполнение приложения
func (a *App) Stop() error {
	return nil
}
