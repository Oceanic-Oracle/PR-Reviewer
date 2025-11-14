package apperrors

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return e.Message
}

var ErrNotFound = &AppError{
	Code:    "NOT_FOUND",
	Message: "resource not found",
}

var ErrTeamExists = &AppError{
	Code:    "TEAM_EXISTS",
	Message: "team_name already exists",
}

var ErrPRExists = &AppError{
	Code: "PR_EXISTS",
	Message: "PR id already exists",
}

var ErrMerged = &AppError{
	Code: "PR_MERGED",
	Message: "cannot reassign on merged PR",
}

var ErrNotAssigned = &AppError{
	Code: "NOT_ASSIGNED",
	Message: "reviewer is not assigned to this PR ",
}

var ErrNoCandidate = &AppError{
	Code: "NO_CANDIDATE",
	Message: "no active replacement candidate in team",
}
