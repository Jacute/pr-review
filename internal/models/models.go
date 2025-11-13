package models

// User представляет участника команды.
// @Description Модель пользователя
type User struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
