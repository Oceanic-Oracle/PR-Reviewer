package team

import (
	"context"
	"pr/internal/repo/user"
)

type UserInterface interface {
	CreateOrUpdateTeamWithUsers(ctx context.Context, team string, users []user.UserModel) (string, []user.UserModel, error)
	GetUsersFromTeam(ctx context.Context, team string) ([]user.UserModel, error)
}
