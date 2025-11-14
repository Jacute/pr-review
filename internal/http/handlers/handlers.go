package handlers

import (
	"context"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
)

type Usecases interface {
	GetTeam(ctx context.Context, name string) ([]*models.Member, error)
	CreateTeam(ctx context.Context, reqDTO *dto.AddTeamRequest) error

	UserSetIsActive(ctx context.Context, reqDTO *dto.SetIsActiveRequest) (*models.User, error)
	GetReviewers(ctx context.Context, userId string) ([]*models.PullRequest, error)

	CreatePR(ctx context.Context, reqDTO *dto.CreatePRRequest) (*models.PullRequest, error)
	MergePR(ctx context.Context, prId string) (*models.PullRequest, error)
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
