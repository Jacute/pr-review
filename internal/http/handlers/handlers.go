package handlers

import (
	"context"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
)

type Usecases interface {
	GetTeam(ctx context.Context, name string) ([]*models.User, error)
	CreateTeam(ctx context.Context, reqDTO *dto.AddTeamRequest) error
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
