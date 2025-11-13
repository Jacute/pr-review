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
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	cfg := config.MustParseConfig()

	db := postgres.New(cfg.DatabaseConfig)
	uc := usecases.New(log, db)
	m := middlewares.New(log)
	h := handlers.New(log, uc)

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT)

	s := server.New(log, cfg.ApplicationConfig, h, m)
	go s.Run()

	<-sign

}
