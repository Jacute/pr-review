package usecases

import (
	"context"
	"errors"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
	"pr-review/internal/repository/postgres"
	"pr-review/internal/utils"
)

var (
	ErrPRNotFound      = errors.New("resource not found")
	ErrPRAlreadyExists = errors.New("PR id already exists")
)

const maxReviewersPerPR = 2

func (uc *Usecases) CreatePR(ctx context.Context, reqDTO *dto.CreatePRRequest) (*models.PullRequest, error) {
	const op = "usecases.CreatePR"
	log := uc.log.With(slog.String("op", op), slog.String("pr_id", reqDTO.Id))

	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		log.Error("error beginning transaction", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error("error rolling back transaction", slog.String("error", rbErr.Error()))
			}
			return
		}
		if cmErr := tx.Commit(ctx); cmErr != nil {
			log.Error("error committing transaction", slog.String("error", cmErr.Error()))
			err = cmErr
		}
	}()

	_, err = uc.db.GetUserById(ctx, reqDTO.AuthorID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			log.Warn("author not found")
			return nil, ErrUserNotFound
		}
		log.Error("error getting user by id", slog.String("error", err.Error()))
		return nil, err
	}

	err = uc.db.CreatePR(ctx, tx, &models.PullRequestShort{
		Id:                reqDTO.Id,
		Title:             reqDTO.Title,
		AuthorId:          reqDTO.AuthorID,
		Status:            models.StatusOpen,
		NeedMoreReviewers: true,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrPRAlreadyExists) {
			log.Warn("PR already exists")
			return nil, ErrPRAlreadyExists
		}
		log.Error("error creating PR", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("PR created successfully")

	members, err := uc.db.GetMembers(ctx, tx, reqDTO.Id, "")
	if err != nil {
		log.Error("error getting members", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("teammates got successfully", slog.Int("teammates_count", len(members)))

	utils.Shuffle(members)

	reviewers := make([]string, 0, maxReviewersPerPR)
	for range maxReviewersPerPR {
		assigneeId, err := uc.db.AssignPRToUser(ctx, tx, reqDTO.Id, members)
		if err != nil {
			log.Error("error assigning PR to user", slog.String("error", err.Error()))
			return nil, err
		}
		if assigneeId != "" {
			reviewers = append(reviewers, assigneeId)
		}
	}

	if len(reviewers) == 2 {
		err = uc.db.UpdatePR(ctx, tx, &models.PullRequestShort{
			Id:                reqDTO.Id,
			NeedMoreReviewers: false,
		})
		if err != nil {
			log.Error("error creating PR", slog.String("error", err.Error()))
			return nil, err
		}
	}

	pr, err := uc.db.GetPRById(ctx, tx, reqDTO.Id)

	log.Debug("PR created successfully")

	return pr, nil
}

func (uc *Usecases) MergePR(ctx context.Context, prId string) (*models.PullRequest, error) {
	const op = "usecases.MergePR"
	log := uc.log.With(slog.String("op", op), slog.String("pr_id", prId))

	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		log.Error("error beginning transaction", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error("error rolling back transaction", slog.String("error", rbErr.Error()))
			}
			return
		}
		if cmErr := tx.Commit(ctx); cmErr != nil {
			log.Error("error committing transaction", slog.String("error", cmErr.Error()))
			err = cmErr
		}
	}()

	pr, err := uc.db.GetPRById(ctx, tx, prId)
	if err != nil {
		log.Error("error getting PR by id", slog.String("error", err.Error()))
		return nil, err
	}
	if pr.Status == models.StatusMerged {
		log.Warn("PR is already merged")
		return pr, nil
	}

	err = uc.db.MergePR(ctx, tx, prId)
	if err != nil {
		if errors.Is(err, postgres.ErrPRNotFound) {
			log.Warn("PR not found")
			return nil, ErrPRNotFound
		}
		log.Error("error merging PR", slog.String("error", err.Error()))
		return nil, err
	}

	pr, err = uc.db.GetPRById(ctx, tx, prId)
	if err != nil {
		log.Error("error getting PR by id", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("PR merged successfully")
	return pr, nil
}
