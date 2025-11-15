package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pr/internal/config"
	"pr/internal/repo"
	pr_postgres "pr/internal/repo/pr/postgres"
	team_postgres "pr/internal/repo/team/postgres"
	user_postgres "pr/internal/repo/user/postgres"
	"pr/internal/server/http"
	"pr/pkg/database"
	"pr/pkg/logger"
	"syscall"
	"time"
)

type Bootstrap struct {
	log *slog.Logger
	cfg *config.Config
}

func (b *Bootstrap) Run() {
	repos, close := b.initRepo()
	defer close()

	srv := http.NewRestApi(&b.cfg.Http, repos, b.log)
	srv.CreateServer()
	defer srv.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func (b *Bootstrap) initRepo() (*repo.Repo, func()) {
	switch b.cfg.Storage.Type {
	case "postgres":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := database.NewConn(ctx, b.cfg.Storage.Url, b.log)
		if err != nil {
			// добавить noop
		}

		teamDb := team_postgres.NewPrPostgres(conn, b.log)
		userDb := user_postgres.NewPrPostgres(conn, b.log)
		prDb := pr_postgres.NewPrPostgres(conn, b.log)

		return repo.NewRepo(teamDb, userDb, prDb), func() {
			conn.Close()
		}
	default:
		// добавить noop
		return nil, func(){}
	}
}

func NewBootstrap() *Bootstrap {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	return &Bootstrap{
		log: log,
		cfg: cfg,
	}
}
