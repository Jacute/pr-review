package usecases

import (
	"context"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
	"pr-review/internal/utils"
)

const maxReviewersPerPR = 2

func (uc *Usecases) CreatePR(ctx context.Context, reqDTO *dto.CreatePRRequest) error {
	const op = "usecases.CreatePR"
	log := uc.log.With(slog.String("op", op), slog.String("pr_id", reqDTO.Id))

	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		log.Error("error beginning transaction", slog.String("error", err.Error()))
		return err
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
	err = uc.db.CreatePR(ctx, tx, &models.PullRequestShort{
		Id:                reqDTO.Id,
		Title:             reqDTO.Title,
		AuthorId:          reqDTO.AuthorID,
		Status:            models.StatusOpen,
		NeedMoreReviewers: true,
	})
	if err != nil {
		log.Error("error creating PR", slog.String("error", err.Error()))
		return err
	}

	log.Debug("PR created successfully")

	members, err := uc.db.GetMembers(ctx, tx, reqDTO.Id, "")
	if err != nil {
		log.Error("error getting members", slog.String("error", err.Error()))
		return err
	}
	log.Debug("teammates got successfully", slog.Int("teammates_count", len(members)))

	utils.Shuffle(members)

	reviewers := make([]string, 0, maxReviewersPerPR)
	for range maxReviewersPerPR {
		assigneeId, err := uc.db.AssignPRToUser(ctx, tx, reqDTO.Id, members)
		if err != nil {
			log.Error("error assigning PR to user", slog.String("error", err.Error()))
			return err
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
			return err
		}
	}

	log.Debug("PR created successfully")

	return nil
}
