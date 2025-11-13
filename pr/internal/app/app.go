package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pr/internal/config"
	"pr/internal/repo"
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
	//repos := b.initRepo()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func (b *Bootstrap) initRepo() (repos interface{}) {
	switch b.cfg.Storage.Type {
	case "postgres":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := database.NewConn(ctx, b.cfg.Storage.Url, b.log)
		if err != nil {
			// добавить noop
		}

		return repo.NewRepo(conn, b.log)
	default:
		// добавить noop
		return nil
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