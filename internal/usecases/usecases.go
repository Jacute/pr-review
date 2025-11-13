package usecases

import (
	"log/slog"
	"pr-review/internal/repository/postgres"
)

type Usecases struct {
	log *slog.Logger
	db  *postgres.Storage
}

func New(log *slog.Logger, db *postgres.Storage) *Usecases {
	return &Usecases{
		log: log,
		db:  db,
	}
}
