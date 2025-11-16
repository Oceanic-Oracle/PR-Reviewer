package dto

import apperrors "pr/internal/app/errors"

// general

type ErrorResponse struct {
	Error apperrors.AppError `json:"error"`
}
