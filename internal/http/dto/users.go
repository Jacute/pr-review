package dto

import (
	"pr-review/internal/models"

	"github.com/google/uuid"
)

var (
	ErrUserIdRequired = Error(
		ErrCodeBadRequest,
		"user_id is required",
	)
	ErrUserIdShouldBeUuid = Error(
		ErrCodeBadRequest,
		"user_id should be uuid",
	)
	ErrIsActiveRequired = Error(
		ErrCodeBadRequest,
		"is_active is required",
	)
)

type SetIsActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive *bool  `json:"is_active"`
}

func (r *SetIsActiveRequest) Validate() *ErrorResponse {
	if r.UserId == "" {
		return ErrUserIdRequired
	}
	if _, err := uuid.Parse(r.UserId); err != nil {
		return ErrUserIdShouldBeUuid
	}
	if r.IsActive == nil {
		return ErrIsActiveRequired
	}
	return nil
}

type SetIsActiveResponse struct {
	User *models.User `json:"user"`
}
