package dto

// Requests
type (
	CreateTeamReq struct {
		TeamName string `json:"team_name"`
		Members  []User `json:"members"`
	}
)

// Responses
type (
	CreateTeamRes struct {
		Team CreateTeamReq `json:"team"`
	}
)

// Errors
type (
	CreateTeamErr struct {
		Error ErrorDetails `json:"error"`
	}
)

// General
type (
	User struct {
		UserId   string `json:"user_id"`
		UserName string `json:"username"`
		IsActive bool   `json:"is_active"`
	}

	ErrorDetails struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)
