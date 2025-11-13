package handlers

import (
	"net/http"
)

func (h *Handlers) AddTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented")
	}
}

// GetTeam godoc
// @Summary Получить команду с участниками
// @Param team_name query string true "Название команды	"
// @Produce json
// @Success 200 {object} dto.GetTeamResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 404 {object} dto.Response "Команда не найдена"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /team/get [get]
// @Tags Teams
func (h *Handlers) GetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented")
	}
}
