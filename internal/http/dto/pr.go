package dto

import (
	"net/url"
	"pr-review/internal/models"
	"strconv"

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
	ErrPageShouldBePositiveInt = Error(
		ErrCodeBadRequest,
		"page should be positive number",
	)
	ErrLimitShouldBePositiveInt = Error(
		ErrCodeBadRequest,
		"limit should be positive number",
	)
)

var (
	ErrCodePRExists               ErrorCode = "PR_EXISTS"
	ErrCodeNoCandidates           ErrorCode = "NO_CANDIDATES"
	ErrCodeCannotReassignMergedPR ErrorCode = "PR_MERGED"
	ErrCodeUserNotReviewerOfPR    ErrorCode = "NOT_ASSIGNED"
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

type StatisticsRequest struct {
	Page  int
	Limit int
}

type StatisticsResponse struct {
	Statistics map[string]int `json:"prs"`
	Count      uint64         `json:"authors_count"`
}

func MapQueryToStatisticsRequest(query url.Values) (*StatisticsRequest, *ErrorResponse) {
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	var err error
	var page, limit int
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			return nil, ErrPageShouldBePositiveInt
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return nil, ErrLimitShouldBePositiveInt
		}
	}

	return &StatisticsRequest{
		Page:  page,
		Limit: limit,
	}, nil
}
