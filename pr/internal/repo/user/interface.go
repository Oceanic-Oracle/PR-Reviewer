package user

import "context"

type UserInterface interface {
	CreateOrUpdateTeamWithUsers(ctx context.Context, team string, users []UserModel) (string, []UserModel, error)
}