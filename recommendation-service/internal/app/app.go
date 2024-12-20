package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"recommendation-service/internal/closer"
	"recommendation-service/internal/handlers"
	"recommendation-service/internal/repository"
	"recommendation-service/internal/service"
	"recommendation-service/log"
	"sync"
)

type IRunner interface {
	Run(ctx context.Context) error
	Stop() error
}

//type IWorker interface {
//	Run(ctx context.Context) error
//	Stop() error
//}

type App struct {
	serverHTTP IRunner
	//worker     IWorker
}

func New() (*App, error) {
	l := log.InitLogger().With(zap.String("app", "recommendation-service"))

	appLoger := log.NewFactory(l)

	db := repository.New()
	closer.Add(db.Close)

	svc := service.New(db)

	httpSrv := handlers.NewServer(appLoger, svc)
	//kafkaWorker, err := NewWorker(svc)
	//if err != nil {
	//	return nil, err
	//}

	return &App{
		serverHTTP: httpSrv,
		//worker:     &kafkaWorker,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	//runConsumer := func(worker IWorker) {
	//	wg.Add(1)
	//	defer wg.Done()
	//
	//	err := worker.Run(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//
	//go runConsumer(a.worker)
	//closer.Add(a.worker.Stop)

	runServ := func(runner IRunner) {
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

	go runServ(a.serverHTTP)
	closer.Add(a.serverHTTP.Stop)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan

	//wg.Wait()
	return nil
}

func (a *App) Stop() error {
	return nil
}
