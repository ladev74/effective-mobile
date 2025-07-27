package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"effmob/internal/api/handlers"
	cconfig "effmob/internal/config"
	llogger "effmob/internal/logger"
	ppostgresClient "effmob/internal/storage/postgresClient"
)

const (
	pathToConfigFile     = "./config/config.env"
	pathToMigrationsFile = "file://./database/migrations"
	shoutdownTime        = 15 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	config, err := cconfig.New(pathToConfigFile)
	if err != nil {
		log.Fatal("failed to initialize config", err)
	}

	logger, err := llogger.New(&config.Logger)
	if err != nil {
		log.Fatal("failed to initialize logger", err)
	}

	postgresClient, err := ppostgresClient.New(ctx, &config.Postgres, logger, pathToMigrationsFile)
	if err != nil {
		log.Fatal("failed to initialize postgres client", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(llogger.MiddlewareLogger(logger, &config.Logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/subscription", handlers.AddSubscriptionHandler(logger, postgresClient))
	router.Delete("/subscription/{id}", handlers.DeleteSubscriptionHandler(logger, postgresClient))
	router.Get("/subscription/{id}", handlers.GetSubscriptionHandler(logger, postgresClient))
	router.Get("/subscription", handlers.ListSubscriptionsHandler(logger, postgresClient))

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.HttpServer.Host, config.HttpServer.Port),
		Handler: router,
	}

	go func() {
		logger.Info("starting http server", zap.String("addr", server.Addr))
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("cannot start http server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	gracefulShutdown(logger, &server, postgresClient)
}

func gracefulShutdown(logger *zap.Logger, srv *http.Server,
	postgresClient ppostgresClient.PostgresClient) {
	logger.Info("received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shoutdownTime)
	defer shutdownCancel()

	logger.Info("shutting down http server")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("cannot shutdown http server", zap.Error(err))
		return
	}

	postgresClient.Close()

	logger.Info("stopping http server", zap.String("addr", srv.Addr))

	logger.Info("application shutdown completed successfully")
}

// TODO: documentation
