package logctx

import (
	"context"

	"github.com/google/uuid"
)

type Ctx struct {
	ReqID  string
	UserID string
	Team   string
	PR     string
}

type keyType int

const key keyType = 0

func WithReqID(ctx context.Context) context.Context {
	reqID := uuid.New().String()

	if c, ok := ctx.Value(key).(Ctx); ok {
		c.ReqID = reqID
		return context.WithValue(ctx, key, c)
	}

	return context.WithValue(ctx, key, Ctx{ReqID: reqID})
}

func WithUserID(ctx context.Context, userID string) context.Context {
	if c, ok := ctx.Value(key).(Ctx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}

	return context.WithValue(ctx, key, Ctx{UserID: userID})
}

func WithTeam(ctx context.Context, team string) context.Context {
	if c, ok := ctx.Value(key).(Ctx); ok {
		c.Team = team
		return context.WithValue(ctx, key, c)
	}

	return context.WithValue(ctx, key, Ctx{Team: team})
}

func WithPR(ctx context.Context, prID string) context.Context {
	if c, ok := ctx.Value(key).(Ctx); ok {
		c.PR = prID
		return context.WithValue(ctx, key, c)
	}

	return context.WithValue(ctx, key, Ctx{PR: prID})
}
