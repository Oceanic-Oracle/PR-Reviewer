package dto

import "time"

// general

type PR struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
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
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type MergePRRequest struct {
	ID string `json:"pull_request_id"`
}

type SwapPRReviewerRequest struct {
	PRID      string `json:"pull_request_id"`
	OldUserID string `json:"old_reviewer_id"`
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
