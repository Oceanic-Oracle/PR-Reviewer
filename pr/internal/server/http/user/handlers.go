package user

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	apperrors "pr/internal/app/errors"
	"pr/internal/dto"
	"pr/internal/repo"
	"pr/internal/repo/user"
	"time"
)

func CreateTeam(repo *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := &dto.CreateTeamReq{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userModelMem := make([]user.UserModel, 0, len(body.Members))
		for _, val := range body.Members {
			userModelMem = append(userModelMem, user.UserModel{
				Id: val.UserId,
				TeamName: body.TeamName,
				IsActive: val.IsActive,
				UserName: val.UserName,
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, _, err := repo.User.CreateOrUpdateTeamWithUsers(ctx, body.TeamName, userModelMem); err != nil {
			if errors.Is(err, apperrors.ErrTeamExists) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.CreateTeamErr{
					Error: dto.ErrorDetails{
						Code: apperrors.ErrTeamExists.Code,
						Message: apperrors.ErrTeamExists.Message,
					},
				})
			} else {
				http.Error(w, "Неизвестная ошибка", http.StatusInternalServerError)
			}
			return
		}

		response := &dto.CreateTeamRes{
			Team: *body,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
