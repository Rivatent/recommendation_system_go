package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"product-service/internal/closer"
	"product-service/internal/handlers"
	"product-service/internal/repository"
	"product-service/internal/service"
	"product-service/log"
	"sync"
)

type IRunner interface {
	Run(ctx context.Context) error
	Stop() error
}

type App struct {
	serverHttp IRunner
}

func New() (*App, error) {
	l := log.InitLogger().With(zap.String("app", "product-service"))

	appLogger := log.NewFactory(l)

	db := repository.New()
	closer.Add(db.Close)

	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	kafkaProd, err := service.NewKafkaProducer(broker, topic)
	if err != nil {
		return nil, err
	}
	closer.Add(kafkaProd.Close)

	svc := service.New(db, kafkaProd)

	httpSrv := handlers.NewServer(appLogger, svc)

	return &App{
		serverHttp: httpSrv,
	}, nil
}

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

func (a *App) Stop() error {
	return nil
}
