package user

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	apperrors "pr/internal/app/errors"
	"pr/internal/dto"
	"pr/internal/logctx"
	"pr/internal/repo"
	errorhandler "pr/internal/server/http/error"
	"time"
)

func SetUserFlag(repos *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())

		log.InfoContext(ctx, "Received SetUserFlag request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var body dto.SetUserStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in request body", slog.Any("error", err))
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		if body.ID == "" {
			log.WarnContext(ctx, "user_id is empty")
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		ctx = logctx.WithUserID(ctx, body.ID)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		user, err := repos.User.SetUserStatus(ctx, body.ID, body.IsActive)
		if err != nil {
			log.ErrorContext(ctx, "Failed to update user status", slog.Any("error", err))
			errorhandler.WriteError(w, err, log)

			return
		}

		response := dto.SetUserStatusResponse{
			User: dto.UserWithTeam{
				UserID:   user.ID,
				Username: user.UserName,
				IsActive: user.IsActive,
				TeamName: user.TeamName,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "User status updated successfully")
	}
}

func GetUserPR(repos *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())
		log.InfoContext(ctx, "Received GetUserPR request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			log.WarnContext(ctx, "Missing user_id query parameter")
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		ctx = logctx.WithUserID(ctx, userID)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		prsModel, err := repos.User.GetPrUser(ctx, userID)
		if err != nil {
			log.ErrorContext(ctx, "Failed to get user PRs", slog.Any("error", err))
			errorhandler.WriteError(w, err, log)

			return
		}

		var prs []dto.PR
		for _, val := range prsModel {
			prs = append(prs, dto.PR{
				ID:       val.ID,
				AuthorID: val.AuthorID,
				Name:     val.Name,
				Status:   val.Status,
			})
		}

		response := dto.GetUserPRResponse{
			UserID: userID,
			PR:     prs,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode GetUserPR response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "User PRs retrieved successfully")
	}
}
