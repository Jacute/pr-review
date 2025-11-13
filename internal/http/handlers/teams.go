package handlers

import (
	"errors"
	"net/http"
	"pr-review/internal/http/dto"
	"pr-review/internal/usecases"

	"github.com/go-chi/render"
)

// AddTeam godoc
// @Summary Создать команду с участниками (создаёт/обновляет пользователей)
// @Param request body dto.AddTeamRequest true "Команда"
// @Produce json
// @Success 201 {object} dto.AddTeamResponse
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос или команда с таким team_name уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /team/add [post]
// @Tags Teams
func (h *Handlers) AddTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.Error(
				dto.ErrCodeBadRequest,
				"Content-Type must be application/json",
			))
			return
		}

		var req dto.AddTeamRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.Error(
				dto.ErrCodeBadRequest,
				err.Error(),
			))
			return
		}
		if err := req.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, err)
			return
		}

		err := h.uc.CreateTeam(r.Context(), &req)
		if err != nil {
			if errors.Is(err, usecases.ErrTeamAlredyExists) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, dto.Error(
					dto.ErrCodeBadRequest,
					err.Error(),
				))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.Error(
				dto.ErrCodeInternal,
				"internal error",
			))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, dto.AddTeamResponse{
			Team: &dto.Team{
				Name:    req.Name,
				Members: req.Members,
			},
		})
	}
}

// GetTeam godoc
// @Summary Получить команду с участниками
// @Param team_name query string true "Название команды"
// @Produce json
// @Success 200 {object} dto.GetTeamResponse
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос"
// @Failure 404 {object} dto.ErrorResponse "Команда не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка"
// @Router /team/get [get]
// @Tags Teams
func (h *Handlers) GetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.Error(
				dto.ErrCodeBadRequest,
				"team_name is required",
			))
			return
		}
		members, err := h.uc.GetTeam(r.Context(), teamName)
		if err != nil {
			if errors.Is(err, usecases.ErrTeamNotFound) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, dto.Error(
					dto.ErrCodeNotFound,
					"team not found",
				))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.Error(
				dto.ErrCodeInternal,
				"internal error",
			))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, dto.GetTeamResponse{
			Name:    teamName,
			Members: members,
		})
	}
}
