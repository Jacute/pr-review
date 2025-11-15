package models

import "time"

type Status string

var (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

type User struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}

type Team struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Member struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type PullRequestShort struct {
	Id                string `json:"pull_request_id"`
	Title             string `json:"pull_request_name"`
	AuthorId          string `json:"author_id"`
	Status            Status `json:"status"`
	NeedMoreReviewers bool   `json:"-"`
}

type PullRequest struct {
	PullRequestShort
	Reviewers []string   `json:"assigned_reviewers"`
	MergedAt  *time.Time `json:"merged_at"`
}
