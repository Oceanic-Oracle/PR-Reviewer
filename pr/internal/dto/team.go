package dto

// general

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

// requests

type CreateTeamRequest struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

// responses

type CreateTeamResponse struct {
	Team Team `json:"team"`
}

type GetTeamResponse struct {
	Team Team `json:"team"`
}
