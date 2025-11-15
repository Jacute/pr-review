package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           pr-review API Docs

// @host      localhost:8080
// @BasePath  /
func initRouter(h Handlers, m Middlewares) http.Handler {
	r := chi.NewRouter()

	middleware.DefaultLogger = m.RequestLogger()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(m.Recoverer)

	r.Get("/team/get", h.GetTeam())
	r.Post("/team/add", h.AddTeam())
	r.Post("/users/setIsActive", h.UserSetIsActive())
	r.Get("/users/getReview", h.GetUserReviews())
	r.Post("/pullRequest/create", h.CreatePR())
	r.Post("/pullRequest/merge", h.MergePR())
	r.Post("/pullRequest/reassign", h.ReassignPR())
	r.Get("/pullRequest/statistics", h.Statistics())

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	return r
}
