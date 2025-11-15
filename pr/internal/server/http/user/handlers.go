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
		ctx := logctx.WithReqId(r.Context())

		log.InfoContext(ctx, "Received SetUserFlag request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var body dto.SetUserStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in request body", slog.Any("error", err))
			errorhandler.WriteError(w, apperrors.ErrBadRequest)
			return
		}

		if body.Id == "" {
			log.WarnContext(ctx, "user_id is empty")
			errorhandler.WriteError(w, apperrors.ErrBadRequest)
			return
		}

		ctx = logctx.WithUserID(ctx, body.Id)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		user, err := repos.User.SetUserStatus(ctx, body.Id, body.IsActive)
		if err != nil {
			log.ErrorContext(ctx, "Failed to update user status", slog.Any("error", err))
			errorhandler.WriteError(w, err)
			return
		}

		response := dto.SetUserStatusResponse{
			User: dto.UserWithTeam{
				UserId:   user.Id,
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
		ctx := logctx.WithReqId(r.Context())
		log.InfoContext(ctx, "Received GetUserPR request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		userId := r.URL.Query().Get("user_id")
		if userId == "" {
			log.WarnContext(ctx, "Missing user_id query parameter")
			errorhandler.WriteError(w, apperrors.ErrBadRequest)
			return
		}

		ctx = logctx.WithUserID(ctx, userId)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		prsModel, err := repos.User.GetPrUser(ctx, userId)
		if err != nil {
			log.ErrorContext(ctx, "Failed to get user PRs", slog.Any("error", err))
			errorhandler.WriteError(w, err)
			return
		}

		var prs []dto.PR
		for _, val := range prsModel {
			prs = append(prs, dto.PR{
				Id:        val.Id,
				AuthorId:  val.AuthorId,
				Name:      val.Name,
				Status:    val.Status,
			})
		}

		response := dto.GetUserPRResponse{
			UserId: userId,
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
