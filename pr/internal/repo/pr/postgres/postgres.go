package pr_postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PrPostgres struct {
	conn *pgxpool.Pool
	log *slog.Logger
}

func NewPrPostgres(conn *pgxpool.Pool, log *slog.Logger) *PrPostgres {
	return &PrPostgres{
		conn: conn,
		log: log,
	}
}