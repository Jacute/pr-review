package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"pr-review/internal/config"
	"time"
)

type Handlers interface {
	CreatePR() http.HandlerFunc
	MergePR() http.HandlerFunc
	ReassignPR() http.HandlerFunc
	GetUserReviews() http.HandlerFunc
	GetTeam() http.HandlerFunc
	AddTeam() http.HandlerFunc
	UserSetIsActive() http.HandlerFunc
	Statistics() http.HandlerFunc
}

type Middlewares interface {
	Recoverer(next http.Handler) http.Handler
	RequestLogger() func(next http.Handler) http.Handler
}

type HTTPServer struct {
	log    *slog.Logger
	cfg    *config.ApplicationConfig
	server *http.Server
}

func New(log *slog.Logger, cfg *config.ApplicationConfig, h Handlers, m Middlewares) *HTTPServer {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      initRouter(h, m),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return &HTTPServer{
		log:    log,
		cfg:    cfg,
		server: srv,
	}
}

func (s *HTTPServer) TestReq(req *http.Request, res http.ResponseWriter) {
	s.server.Handler.ServeHTTP(res, req)
}

func (a *HTTPServer) Run() {
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic("HTTP server failed to start: " + err.Error())
	}
}

func (a *HTTPServer) Stop() {
	const op = "server.Stop"
	log := a.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		log.Error("failed to graceful shutdown http server, trying hard close", slog.String("error", err.Error()))
		err = a.server.Close()
		if err != nil {
			log.Error("failed to hard close server", slog.String("error", err.Error()))
		}
	}
}
