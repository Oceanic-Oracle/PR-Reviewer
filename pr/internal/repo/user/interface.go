package user

import (
	"context"
	"pr/internal/repo/pr"
)

type UserInterface interface {
	SetUserStatus(ctx context.Context, id string, isActive bool) (UserModel, error)
	GetPrUser(ctx context.Context, userId string) ([]pr.PRModel, error)
}