package dto

import (
	"pr-review/internal/models"

	"github.com/google/uuid"
)

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
	ErrOldReviewerIdRequired = Error(
		ErrCodeBadRequest,
		"old_reviewer_id is required",
	)
	ErrOldReviewerIdShouldBeUuid = Error(
		ErrCodeBadRequest,
		"old_reviewer_id should be uuid",
	)
)

var (
	ErrCodePRExists     ErrorCode = "PR_EXISTS"
	ErrCodeNoCandidates ErrorCode = "NO_CANDIDATES"
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
	PR *models.PullRequest `json:"pr"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (r *MergePRRequest) Validate() *ErrorResponse {
	if r.PullRequestID == "" {
		return ErrPRIdRequired
	}
	if _, err := uuid.Parse(r.PullRequestID); err != nil {
		return ErrPRIdShouldBeUuid
	}
	return nil
}

type MergePRResponse struct {
	PR *models.PullRequest `json:"pr"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

func (r *ReassignPRRequest) Validate() *ErrorResponse {
	if r.PullRequestID == "" {
		return ErrPRIdRequired
	}
	if _, err := uuid.Parse(r.PullRequestID); err != nil {
		return ErrPRIdShouldBeUuid
	}
	if r.OldReviewerID == "" {
		return ErrOldReviewerIdRequired
	}
	if _, err := uuid.Parse(r.OldReviewerID); err != nil {
		return ErrOldReviewerIdShouldBeUuid
	}
	return nil
}

type ReassignPRResponse struct {
	PR         *models.PullRequest `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}
