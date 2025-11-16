package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pr/internal/config"
	"pr/internal/repo"
	prpostgres "pr/internal/repo/pr/postgres"
	teampostgres "pr/internal/repo/team/postgres"
	userpostgres "pr/internal/repo/user/postgres"
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
	repos, closeDB := b.initRepo()
	defer closeDB()

	srv := http.NewRestAPI(&b.cfg.HTTP, repos, b.log)
	srv.CreateServer()
	defer srv.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	b.log.Info("shutting down...")
}

func (b *Bootstrap) initRepo() (*repo.Repo, func()) {
	switch b.cfg.Storage.Type {
	case "postgres":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		conn, err := database.NewConn(ctx, b.cfg.Storage.URL, b.log)

		cancel()

		if err != nil {
			// добавить noop
			b.log.Error("failed to connect to database", "error", err)
			os.Exit(1)

			return nil, func() {}
		}

		teamDB := teampostgres.NewPrPostgres(conn, b.log)
		userDB := userpostgres.NewPrPostgres(conn, b.log)
		prDB := prpostgres.NewPrPostgres(conn, b.log)

		return repo.NewRepo(teamDB, userDB, prDB), func() {
			conn.Close()
		}
	default:
		// добавить noop
		b.log.Error("unsupported storage type", "type", b.cfg.Storage.Type)
		os.Exit(1)

		return nil, func() {}
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
