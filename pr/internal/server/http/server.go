package http

import (
	"log/slog"
	"net/http"
	"pr/internal/config"
	"pr/internal/repo"
	"pr/internal/server/http/pr"
	"pr/internal/server/http/team"
	"pr/internal/server/http/user"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	log  *slog.Logger
	cfg  *config.HTTP
	repo *repo.Repo

	// Do not touch
	srv *http.Server
}

func (s *Server) CreateServer() {
	router := chi.NewRouter()

	router.Post("/team/add", team.CreateTeam(s.repo, s.log))
	router.Get("/team/get", team.GetTeam(s.repo, s.log))
	router.Post("/users/setIsActive", user.SetUserFlag(s.repo, s.log))
	router.Post("/pullRequest/create", pr.CreatePR(s.repo, s.log))
	router.Post("/pullRequest/merge", pr.MergePR(s.repo, s.log))
	router.Post("/pullRequest/reassign", pr.SwapPRReviewer(s.repo, s.log))
	router.Get("/users/getReview", user.GetUserPR(s.repo, s.log))

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})

	idleTimeout, _ := strconv.Atoi(s.cfg.IdleTimeout)
	timeout, _ := strconv.Atoi(s.cfg.Timeout)
	srv := &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      corsHandler,
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

func NewRestAPI(cfg *config.HTTP, repo *repo.Repo, log *slog.Logger) *Server {
	return &Server{
		cfg:  cfg,
		repo: repo,
		log:  log,
	}
}
