package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(ctx context.Context, url string, log *slog.Logger) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Error("failed to connection to postgres", "err", err)

		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		log.Error("failed to ping database",
			"err", err,
			"url", url)

		return nil, err
	}

	return conn, nil
}
