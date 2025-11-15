package logctx

import (
	"context"

	"github.com/google/uuid"
)

type Ctx struct {
	ReqId  string
	UserId string
	Team   string
	PR     string
}

type keyType int

const key keyType = 0

func WithReqId(ctx context.Context) context.Context {
	reqId := uuid.New().String()

	if c, ok := ctx.Value(key).(Ctx); ok {
		c.ReqId = reqId
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, Ctx{ReqId: reqId})
}

func WithUserID(ctx context.Context, userId string) context.Context {
	if c, ok := ctx.Value(key).(Ctx); ok {
		c.UserId = userId
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, Ctx{UserId: userId})
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