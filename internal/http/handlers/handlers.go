package handlers

import "log/slog"

type Usecases interface {
}

type Handlers struct {
	log *slog.Logger
	uc  Usecases
}

func New(log *slog.Logger, uc Usecases) *Handlers {
	return &Handlers{
		log: log,
		uc:  uc,
	}
}
