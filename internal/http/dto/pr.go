package dto

import "github.com/google/uuid"

var (
	ErrPRIdRequired = Error(
		ErrCodeBadRequest,
		"pull_request_id is required",
	)
	ErrPRIdShouldBeUuid = Error(
		ErrCodeBadRequest,
		"pull_request_id should be uuid",
	)
	ErrPRTitleRequired = Error(
		ErrCodeBadRequest,
		"pull_request_name is required",
	)
	ErrAuthorIdRequired = Error(
		ErrCodeBadRequest,
		"author_id is required",
	)
	ErrAuthorIdShouldBeUuid = Error(
		ErrCodeBadRequest,
		"author_id should be uuid",
	)
)

type CreatePRRequest struct {
	Id       string `json:"pull_request_id"`
	Title    string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func (r *CreatePRRequest) Validate() *ErrorResponse {
	if r.Id == "" {
		return ErrPRIdRequired
	}
	if _, err := uuid.Parse(r.Id); err != nil {
		return ErrPRIdShouldBeUuid
	}
	if r.Title == "" {
		return ErrPRTitleRequired
	}
	if r.AuthorID == "" {
		return ErrAuthorIdRequired
	}
	if _, err := uuid.Parse(r.AuthorID); err != nil {
		return ErrAuthorIdShouldBeUuid
	}
	return nil
}

type CreatePRResponse struct {
	PullRequestID string `json:"pull_request_id"`
}
