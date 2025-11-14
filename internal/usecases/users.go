package usecases

import (
	"context"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
)

func (uc *Usecases) UserSetIsActive(ctx context.Context, reqDTO *dto.SetIsActiveRequest) (*models.User, error) {
	const op = "usecases.UserSetIsActive"
	log := uc.log.With(slog.String("op", op), slog.String("user_id", reqDTO.UserId))

	err := uc.db.UserSetIsActive(ctx, reqDTO)
	if err != nil {
		log.Error("error setting isActive", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("set is_active successfully")

	user, err := uc.db.GetUserById(ctx, reqDTO.UserId)
	if err != nil {
		log.Error("error getting user by id", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("user got successfully", slog.Any("user", user))

	return user, nil
}

func (uc *Usecases) GetReviewers(ctx context.Context, userId string) ([]*models.PullRequest, error) {
	const op = "usecases.GetReviewers"
	log := uc.log.With(slog.String("op", op), slog.String("user_id", userId))

	reviewers, err := uc.db.GetPRsByUserId(ctx, userId)
	if err != nil {
		log.Error("error getting reviewers by user id", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("reviewers got successfully", slog.Int("reviewers_count", len(reviewers)))

	return reviewers, nil
}
