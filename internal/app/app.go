package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/M1r0-dev/Subscription-Aggregator/config"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/repo/persistence"
	subscriptionservice "github.com/M1r0-dev/Subscription-Aggregator/internal/usecase/subscriptionService"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/httpserver"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/logger"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/postgres"
)

func Run(cfg *config.Config) {
	//Logger
	l := logger.New(cfg.Log.Level)

	//Postgres
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	//Usecase
	SubscriptionUsecase := subscriptionservice.New(
		persistence.New(pg),
	)

	//http server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, SubscriptionUsecase, l)

	httpServer.Start()


	//Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}
	
	//Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
