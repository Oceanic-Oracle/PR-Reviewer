package errorhandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	apperrors "pr/internal/app/errors"
	"pr/internal/dto"
)

func WriteError(w http.ResponseWriter, err error, log *slog.Logger) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		resp := dto.ErrorResponse{
			Error: *appErr,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Status)

		if encodeErr := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode error response", slog.Any("original_error", err), slog.Any("encode_error", encodeErr))
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	if encodeErr := json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: *apperrors.ErrInternalServer,
	}); err != nil {
		log.Error("Failed to encode internal error response",
			slog.Any("original_error", err), slog.Any("encode_error", encodeErr))
	}
}
