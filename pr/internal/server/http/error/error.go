package errorhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	apperrors "pr/internal/app/errors"
	"pr/internal/dto"
)

func WriteError(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		resp := dto.ErrorResponse{
			Error: *appErr,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: *apperrors.ErrInternalServer,
	})
}
