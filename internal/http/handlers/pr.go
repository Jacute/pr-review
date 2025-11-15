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
// @Router /pullRequest/create [post]
// @Tags PullRequests
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

// MergePR godoc
// @Summary Пометить PR как MERGED (идемпотентная операция)
// @Param request body dto.MergePRRequest true "PR id"
// @Produce json
// @Success 200 {object} dto.MergePRResponse "PR в состоянии MERGED"
// @Failure 404 {object} dto.ErrorResponse "PR не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /pullRequest/merge [post]
// @Tags PullRequests
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

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, dto.MergePRResponse{
			PR: pr,
		})
	}
}

// ReassignPR godoc
// @Summary Переназначить конкретного ревьювера на другого из его команды
// @Param request body dto.ReassignPRRequest true "PR id & old reviewer id"
// @Produce json
// @Success 200 {object} dto.ReassignPRResponse
// @Failure 404 {object} dto.ErrorResponse "PR или пользователь не найден"
// @Failure 409 {object} dto.ErrorResponse "Нельзя менять после MERGED"
// @Failure 409 {object} dto.ErrorResponse "Пользователь не был назначен ревьювером"
// @Failure 409 {object} dto.ErrorResponse "Нет доступных кандидатов"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /pullRequest/reassign [post]
// @Tags PullRequests
func (h *Handlers) ReassignPR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrContentTypeNotJson)
			return
		}

		var req dto.ReassignPRRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrInvalidBody)
			return
		}

		pr, replacedBy, err := h.uc.ReassignPR(r.Context(), &req)
		if err != nil {
			if errors.Is(err, usecases.ErrPRNotFound) || errors.Is(err, usecases.ErrUserNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, dto.Error(dto.ErrCodeNotFound, err.Error()))
				return
			}
			if errors.Is(err, usecases.ErrNoCandidatesToAssign) {
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, dto.Error(dto.ErrCodeNoCandidates, err.Error()))
				return
			}
			if errors.Is(err, usecases.ErrPRMerged) {
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, dto.Error(dto.ErrCodeCannotReassignMergedPR, err.Error()))
				return
			}
			if errors.Is(err, usecases.ErrUserNotReviewerOfPR) {
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, dto.Error(dto.ErrCodeUserNotReviewerOfPR, err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, dto.ReassignPRResponse{
			PR:         pr,
			ReplacedBy: replacedBy,
		})
	}
}
