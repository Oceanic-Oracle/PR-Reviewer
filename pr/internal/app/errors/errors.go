// internal/apperrors/apperrors.go
package apperrors

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrTeamExists = &AppError{
		Code:    "TEAM_EXISTS",
		Message: "team_name already exists",
		Status:  http.StatusBadRequest,
	}

	ErrNotFound = &AppError{
		Code:    "NOT_FOUND",
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}

	ErrInternalServer = &AppError{
		Code:    "INTERNAL",
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}

	ErrBadRequest = &AppError{
		Code:    "BAD_REQUEST",
		Message: "bad request",
		Status:  http.StatusBadRequest,
	}

	ErrPRExists = &AppError{
		Code:    "PR_EXISTS",
		Message: "PR id already exists",
		Status:  http.StatusConflict,
	}

	ErrMerged = &AppError{
		Code:    "PR_MERGED",
		Message: "cannot reassign on merged PR",
		Status:  http.StatusConflict,
	}

	ErrNotAssigned = &AppError{
		Code:    "NOT_ASSIGNED",
		Message: "reviewer is not assigned to this PR",
		Status:  http.StatusConflict,
	}

	ErrNoCandidate = &AppError{
		Code:    "NO_CANDIDATE",
		Message: "no active replacement candidate in team",
		Status:  http.StatusConflict,
	}
)
