package e2e

import (
	"log/slog"
	"os"
	"pr-review/internal/config"
	"pr-review/internal/http/handlers"
	"pr-review/internal/http/middlewares"
	"pr-review/internal/http/server"
	"pr-review/internal/repository/postgres"
	"pr-review/internal/usecases"

	"github.com/joho/godotenv"
)

const envconfigFilename = ".env.test"

type Suite struct {
	srv *server.HTTPServer
	db  *postgres.Storage
}

func NewSuite() *Suite {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	godotenv.Load(envconfigFilename)
	cfg := config.MustParseConfig()

	db := postgres.New(cfg.DatabaseConfig)
	uc := usecases.New(log, db)
	h := handlers.New(log, uc)
	m := middlewares.New(log)

	srv := server.New(log, cfg.ApplicationConfig, h, m)

	return &Suite{
		srv: srv,
		db:  db,
	}
}

func (s *Suite) Start() {
	go s.srv.Run()
}

func (s *Suite) Stop() {
	s.srv.Stop()
}
