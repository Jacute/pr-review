package handlers

import (
	"errors"
	"net/http"
	"pr-review/internal/http/dto"
	"pr-review/internal/usecases"

	"github.com/go-chi/render"
)

// CreatePR godoc
// @Summary Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// @Param request body dto.CreatePRRequest true "PR"
// @Produce json
// @Success 201 {object} dto.CreatePRResponse
// @Failure 404 {object} dto.ErrorResponse "Автор не найден"
// @Failure 409 {object} dto.ErrorResponse "PR уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /team/add [post]
// @Tags Teams
func (h *Handlers) CreatePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrContentTypeNotJson)
			return
		}

		var req dto.CreatePRRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrInvalidBody)
			return
		}
		if req.Validate() != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, req.Validate())
			return
		}

		pr, err := h.uc.CreatePR(r.Context(), &req)
		if err != nil {
			if errors.Is(err, usecases.ErrUserNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, dto.Error(dto.ErrCodeNotFound, err.Error()))
				return
			}
			if errors.Is(err, usecases.ErrPRAlreadyExists) {
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, dto.Error(dto.ErrCodePRExists, err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, dto.CreatePRResponse{
			PR: pr,
		})
	}
}

func (h *Handlers) MergePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrContentTypeNotJson)
			return
		}

		var req dto.MergePRRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrInvalidBody)
			return
		}
		if req.Validate() != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, req.Validate())
			return
		}

		pr, err := h.uc.MergePR(r.Context(), req.PullRequestID)
		if err != nil {
			if errors.Is(err, usecases.ErrPRNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, dto.Error(dto.ErrCodeNotFound, err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, dto.MergePRResponse{
			PR: pr,
		})
	}
}

func (h *Handlers) ReassignPR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented")
	}
}
