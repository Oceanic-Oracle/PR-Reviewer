package http

import (
	"log/slog"
	"net/http"
	"pr/internal/config"
	"pr/internal/repo"
	"pr/internal/server/http/user"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	log  *slog.Logger
	cfg  *config.Http
	repo *repo.Repo

	//Do not touch
	srv *http.Server
}

func (s *Server) CreateServer() {
	router := chi.NewRouter()

	router.Post("/team/add", user.CreateTeam(s.repo, s.log))

	idleTimeout, _ := strconv.Atoi(s.cfg.IdleTimeout)
	timeout, _ := strconv.Atoi(s.cfg.Timeout)
	srv := &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      router,
		ReadTimeout:  time.Duration(timeout) * time.Second,
		WriteTimeout: time.Duration(timeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	go func() {
		s.log.Info("HTTP server starting", slog.String("addr", s.cfg.Addr))
		if err := srv.ListenAndServe(); err != nil {
			s.log.Error("HTTP server failed", slog.Any("error", err))
			return
		}
	}()

	s.srv = srv
}

func (s *Server) Close() error {
	return s.srv.Close()
}

func NewRestApi(cfg *config.Http, repo *repo.Repo, log *slog.Logger) *Server {
	return &Server{
		cfg:  cfg,
		repo: repo,
		log:  log,
	}
}