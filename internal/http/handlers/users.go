package handlers

import (
	"errors"
	"net/http"
	"pr-review/internal/http/dto"
	"pr-review/internal/usecases"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// UserSetIsActive godoc
// @Summary Установить флаг активности пользователя
// @Param request body dto.SetIsActiveRequest true "Установка пользователя активным/неактивным"
// @Produce json
// @Success 200 {object} dto.SetIsActiveResponse
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос"
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /users/setIsActive [post]
// @Tags Users
func (h *Handlers) UserSetIsActive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Header.Set(w.Header(), "Content-Type", "application/json")
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrContentTypeNotJson)
			return
		}

		var req dto.SetIsActiveRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrInvalidBody)
			return
		}
		if err := req.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		user, err := h.uc.UserSetIsActive(r.Context(), &req)
		if err != nil {
			if errors.Is(err, usecases.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, dto.ErrUserNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, dto.SetIsActiveResponse{
			User: user,
		})
	}
}

// GetUserReviews godoc
// @Summary Получить PR'ы, где пользователь назначен ревьювером
// @Param user_id query string true "Идентификатор пользователя"
// @Produce json
// @Success 200 {object} dto.GetReviewResponse
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос"
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /users/getReview [get]
// @Tags Users
func (h *Handlers) GetUserReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Header.Set(w.Header(), "Content-Type", "application/json")
		userId := r.URL.Query().Get("user_id")
		if userId == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrUserIdRequired)
			return
		}
		if _, err := uuid.Parse(userId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrUserIdShouldBeUuid)
			return
		}

		reviews, err := h.uc.GetPRs(r.Context(), userId)
		if err != nil {
			if errors.Is(err, usecases.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, dto.ErrUserNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, dto.GetReviewResponse{
			UserId:       userId,
			PullRequests: reviews,
		})
	}
}
