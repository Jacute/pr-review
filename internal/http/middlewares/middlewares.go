package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type Middlewares struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Middlewares {
	return &Middlewares{log: log}
}

func (m *Middlewares) Recoverer(next http.Handler) http.Handler {
	const op = "middlewares.Recoverer"
	log := m.log.With(slog.String("op", op))

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if reqID, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
					log = log.With(slog.String("request_id", reqID))
				}
				log.Error(
					"recovered from panic",
					slog.Any("rvr", rvr),
					slog.Any("stack", string(debug.Stack())),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (m *Middlewares) RequestLogger() func(next http.Handler) http.Handler {
	const op = "middlewares.RequestLogger"
	log := m.log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				reqId := middleware.GetReqID(r.Context())
				log.Info(
					"request completed",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("time", time.Since(t1)),
					slog.String("request_id", reqId),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
