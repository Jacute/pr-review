package dto

import (
	"pr-review/internal/models"

	"github.com/google/uuid"
)

var (
	ErrCodeTeamExists ErrorCode = "TEAM_EXISTS"
)

var (
	ErrTeamAlreadyExists = Error(
		ErrCodeTeamExists,
		"team already exists",
	)
	ErrTeamNameRequired = Error(
		ErrCodeBadRequest,
		"team_name is required",
	)
	ErrTeamNotFound = Error(
		ErrCodeNotFound,
		"team not found",
	)
)

type GetTeamResponse struct {
	Name    string           `json:"team_name"`
	Members []*models.Member `json:"members"`
}

type AddTeamRequest struct {
	Name    string           `json:"team_name"`
	Members []*models.Member `json:"members"`
}

func (r *AddTeamRequest) Validate() *ErrorResponse {
	if r.Name == "" {
		return &ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "team_name is required",
		}
	}
	if len(r.Name) > 255 {
		return &ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "team_name is too long",
		}
	}
	if len(r.Members) > 300 {
		return &ErrorResponse{
			Code:    ErrCodeBadRequest,
			Message: "too many members",
		}
	}
	for _, m := range r.Members {
		if m.Id == "" {
			return &ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "user_id is required",
			}
		}
		if _, err := uuid.Parse(m.Id); err != nil {
			return &ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "user_id should be uuid",
			}
		}
		if m.Username == "" {
			return &ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "username is required",
			}
		}
		if len(m.Username) > 255 {
			return &ErrorResponse{
				Code:    ErrCodeBadRequest,
				Message: "team_name is too long",
			}
		}
	}
	return nil
}

type Team struct {
	Name    string           `json:"team_name"`
	Members []*models.Member `json:"members"`
}

type AddTeamResponse struct {
	Team *Team `json:"team"`
}
