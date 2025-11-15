package dto

// general

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type UserWithTeam struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}

// requests

type SetUserStatusRequest struct {
	Id       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// responses

type SetUserStatusResponse struct {
	User UserWithTeam `json:"user"`
}

type GetUserPRResponse struct {
	UserId string `json:"user_id"`
	PR     []PR   `json:"pull_requests"`
}
