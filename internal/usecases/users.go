package usecases

import (
	"context"
	"errors"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
	"pr-review/internal/repository/postgres"
)

var (
	ErrUserNotFound = errors.New("resourse not found")
)

func (uc *Usecases) UserSetIsActive(ctx context.Context, reqDTO *dto.SetIsActiveRequest) (*models.User, error) {
	const op = "usecases.UserSetIsActive"
	log := uc.log.With(slog.String("op", op), slog.String("user_id", reqDTO.UserId))

	err := uc.db.UserSetIsActive(ctx, reqDTO)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			log.Error("user not found", slog.String("user_id", reqDTO.UserId))
			return nil, ErrUserNotFound
		}
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

func (uc *Usecases) GetPRs(ctx context.Context, userId string) ([]*models.PullRequest, error) {
	const op = "usecases.GetReviewers"
	log := uc.log.With(slog.String("op", op), slog.String("user_id", userId))

	prs, err := uc.db.GetPRsByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			log.Error("user not found", slog.String("user_id", userId))
			return nil, ErrUserNotFound
		}
		log.Error("error getting prs by user id", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("prs got successfully", slog.Int("ors_count", len(prs)))

	return prs, nil
}
