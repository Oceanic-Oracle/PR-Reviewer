package repo

import "log/slog"

type Repo struct {
	log *slog.Logger
	conn interface{}
}

func NewRepo(conn interface{}, log *slog.Logger) *Repo {
	return &Repo{
		conn: conn,
		log: log,
	}
}
