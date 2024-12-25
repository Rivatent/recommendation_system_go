package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"user-service/internal/closer"
	"user-service/internal/handlers"
	"user-service/internal/monitoring"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/internal/validator"
	"user-service/log"
)

// IRunner определяет интерфейс для запуска и остановки компонентов приложения.
type IRunner interface {
	Run(ctx context.Context) error
	Stop() error
}

// App содержит все необходимые компоненты для работы микросервиса.
type App struct {
	serverHttp IRunner
}

// New создает новое приложение, инициализируя логгер, валидатор, базу данных, продюсер Kafka и сервер.
func New() (*App, error) {
	l := log.InitLogger().With(zap.String("app", "user-service"))

	appLogger := log.NewFactory(l)
	validator.InitValidator()
	monitoring.InitMetrics()

	db := repository.New()
	closer.Add(db.Close)

	kafkaProd := service.NewKafkaProducer(os.Getenv("KAFKA_BROKER"), os.Getenv("KAFKA_TOPIC_NEW"), os.Getenv("KAFKA_TOPIC_UPDATE"))

	closer.Add(kafkaProd.Close)

	svc := service.New(db, kafkaProd)

	httpSrv := handlers.NewServer(appLogger, svc)

	return &App{
		serverHttp: httpSrv,
	}, nil
}

// Run запускает приложение, обрабатывая прерывание с помощью сигнала и управляя жизненным циклом сервера.
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

	go run(a.serverHttp)
	closer.Add(a.serverHttp.Stop)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan

	return nil
}

// Stop останавливает выполнение приложения
func (a *App) Stop() error {
	return nil
}
