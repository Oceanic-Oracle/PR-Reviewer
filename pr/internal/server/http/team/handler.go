package team

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	apperrors "pr/internal/app/errors"
	"pr/internal/dto"
	"pr/internal/logctx"
	"pr/internal/repo"
	"pr/internal/repo/user"
	errorhandler "pr/internal/server/http/error"
)

func CreateTeam(repo *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqId(r.Context())
		log.InfoContext(ctx, "Received CreateTeam request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		body := &dto.CreateTeamRequest{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in CreateTeam request", slog.Any("error", err))
			errorhandler.WriteError(w, apperrors.ErrBadRequest)
			return
		}

		if body.TeamName == "" {
			log.WarnContext(ctx, "Missing team_name in CreateTeam request")
			errorhandler.WriteError(w, apperrors.ErrBadRequest)
			return
		}

		ctx = logctx.WithTeam(ctx, body.TeamName)

		userModelMem := make([]user.UserModel, 0, len(body.Members))
		for _, val := range body.Members {
			if val.UserId == "" {
				log.WarnContext(ctx, "Missing user_id in team member", slog.String("username", val.Username))
				errorhandler.WriteError(w, apperrors.ErrBadRequest)
				return
			}
			userModelMem = append(userModelMem, user.UserModel{
				Id:       val.UserId,
				TeamName: body.TeamName,
				IsActive: val.IsActive,
				UserName: val.Username,
			})
			ctx = logctx.WithUserID(ctx, val.UserId)
		}

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_, users, err := repo.Team.CreateOrUpdateTeamWithUsers(ctx, body.TeamName, userModelMem)
		if err != nil {
			log.ErrorContext(ctx, "Failed to create team", slog.Any("error", err))
			errorhandler.WriteError(w, err)
			return
		}

		var mem []dto.User
		for _, val := range users {
			temp := dto.User{
				UserId:   val.Id,
				Username: val.UserName,
				IsActive: val.IsActive,
			}
			mem = append(mem, temp)
		}

		response := &dto.CreateTeamResponse{
			Team: dto.Team{
				TeamName: body.TeamName,
				Members:  mem,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode CreateTeam response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "Team created successfully")
	}
}

func GetTeam(repos *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqId(r.Context())
		log.InfoContext(ctx, "Received GetTeam request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			log.WarnContext(ctx, "Missing team_name query parameter")
			errorhandler.WriteError(w, apperrors.ErrBadRequest)
			return
		}

		ctx = logctx.WithTeam(ctx, teamName)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		mems, err := repos.Team.GetUsersFromTeam(ctx, teamName)
		if err != nil {
			log.ErrorContext(ctx, "Failed to get team", slog.Any("error", err))
			errorhandler.WriteError(w, err)
			return
		}

		var memResp []dto.User
		for _, val := range mems {
			memResp = append(memResp, dto.User{
				UserId:   val.Id,
				Username: val.UserName,
				IsActive: val.IsActive,
			})
		}

		response := dto.GetTeamResponse{
			Team: dto.Team{
				TeamName: teamName,
				Members:  memResp,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode GetTeam response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "Team retrieved successfully")
	}
}