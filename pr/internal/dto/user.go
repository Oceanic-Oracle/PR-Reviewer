package dto

// general

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type UserWithTeam struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}

// requests

type SetUserStatusRequest struct {
	ID       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// responses

type SetUserStatusResponse struct {
	User UserWithTeam `json:"user"`
}

type GetUserPRResponse struct {
	UserID string `json:"user_id"`
	PR     []PR   `json:"pull_requests"`
}
