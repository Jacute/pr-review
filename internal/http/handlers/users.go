package handlers

import (
	"net/http"
	"pr-review/internal/http/dto"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (h *Handlers) UserSetIsActive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, dto.SetIsActiveResponse{
			User: user,
		})
	}
}

func (h *Handlers) GetUserReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		reviews, err := h.uc.GetReviewers(r.Context(), userId)
		if err != nil {
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
