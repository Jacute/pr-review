package dto

var (
	ErrCodeTeamExists ErrorCode = "TEAM_EXISTS"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
)

type Member struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type GetTeamResponse struct {
	Name    string `json:"team_name"`
	Members []*Member
}
