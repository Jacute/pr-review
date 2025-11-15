package usecases

import (
	"context"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)

	UnassignPRsFromUser(ctx context.Context, tx pgx.Tx, userId string) ([]string, error)
	UnassignPRFromUser(ctx context.Context, tx pgx.Tx, prId string, userId string) error
	GetMembers(ctx context.Context, tx pgx.Tx, prId string) ([]*models.Member, error)
	AssignPRToUser(ctx context.Context, tx pgx.Tx, prId string, members []*models.Member) (string, error)
	SetNeedMoreReviewers(ctx context.Context, tx pgx.Tx, prId string) error
	GetPRsByUserId(ctx context.Context, id string) ([]*models.PullRequest, error)
	CreatePR(ctx context.Context, tx pgx.Tx, pr *models.PullRequestShort) error
	UpdatePR(ctx context.Context, tx pgx.Tx, pr *models.PullRequestShort) error
	GetPRById(ctx context.Context, tx pgx.Tx, id string) (*models.PullRequest, error)
	MergePR(ctx context.Context, tx pgx.Tx, id string) error

	TeamExists(ctx context.Context, name string) (bool, error)
	CreateTeam(ctx context.Context, tx pgx.Tx, team *models.Team) error

	GetTeamMembers(ctx context.Context, name string) ([]*models.Member, error)
	AddOrUpdateTeamMembers(ctx context.Context, tx pgx.Tx, teamId string, members []*models.Member) error
	UserSetIsActive(ctx context.Context, reqDTO *dto.SetIsActiveRequest) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
}

type Usecases struct {
	log *slog.Logger
	db  Storage
}

func New(log *slog.Logger, db Storage) *Usecases {
	return &Usecases{
		log: log,
		db:  db,
	}
}
