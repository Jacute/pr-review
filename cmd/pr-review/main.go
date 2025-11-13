package main

import (
	"log/slog"
	"os"
	"os/signal"
	"pr-review/internal/config"
	"pr-review/internal/http/handlers"
	"pr-review/internal/http/middlewares"
	"pr-review/internal/http/server"
	"pr-review/internal/repository/postgres"
	"pr-review/internal/usecases"
	"syscall"
)

func main() {
	cfg := config.MustParseConfig()

	var logLevel slog.Level
	if cfg.Env == "prod" {
		logLevel = slog.LevelInfo
	} else {
		logLevel = slog.LevelDebug
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	db := postgres.New(cfg.DatabaseConfig)
	uc := usecases.New(log, db)
	m := middlewares.New(log)
	h := handlers.New(log, uc)

	s := server.New(log, cfg.ApplicationConfig, h, m)

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGTERM, syscall.SIGINT)

	log.Info("starting http server")
	go s.Run()

	sign := <-signCh
	log.Info("stopping http server", slog.String("signal", sign.String()))
	db.Stop()
	s.Stop()
}
