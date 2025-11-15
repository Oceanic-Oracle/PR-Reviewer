package dto

import "time"

// general

type PR struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
	Status   string `json:"status"`
}

type PRWithReviewers struct {
	PR
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type PRWithReviewersAndMerge struct {
	PR
	AssignedReviewers []string   `json:"assigned_reviewers"`
	MergedAt          *time.Time `json:"mergedAt"`
}

// requests

type CreatePRRequest struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
}

type MergePRRequest struct {
	Id string `json:"pull_request_id"`
}

type SwapPRReviewerRequest struct {
	PRId      string `json:"pull_request_id"`
	OldUserId string `json:"old_reviewer_id"`
}

// responses

type CreatePRResponse struct {
	PR PRWithReviewers `json:"pr"`
}

type MergePRResponse struct {
	PR PRWithReviewersAndMerge `json:"pr"`
}

type SwapPRReviewerResponse struct {
	PR         PRWithReviewers `json:"pr"`
	ReplacedBy string          `json:"replaced_by"`
}
