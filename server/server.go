package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ohhfishal/marvel-snap-archetype/assets"
	"github.com/ohhfishal/marvel-snap-archetype/stats"
	"github.com/ohhfishal/marvel-snap-archetype/templates/components"
	"github.com/ohhfishal/marvel-snap-archetype/templates/pages"
)

type CMD struct {
	Port           string               `default:"8080" env:"PORT" short:"P" help:"Port to serve on"`
	Host           string               `default:"localhost" env:"HOST" short:"H" help:"Address to serve from"`
	RequestTimeout time.Duration        `default:"30s" help:"How long to keep requests alive"`
	Stats          stats.ServiceOptions `embed:""`
}

func (config *CMD) Run(ctx context.Context, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	statsService, err := stats.NewService(logger, config.Stats)
	if err != nil {
		return fmt.Errorf("creating stats service: %w", err)
	}

	r := chi.NewRouter()

	r.Use(loggingMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(config.RequestTimeout))

	r.Mount(
		"/assets",
		http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets))),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home(pages.HomeProps{
			Tournaments: statsService.Tournaments(),
		}).Render(r.Context(), w)
	})

	r.Get("/components/dashboard", func(w http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		archetypes, err := statsService.GetArchetypes(r.Context(), tid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
			logger.Warn("error getting dashboard", "error", err, "tid", tid)
			return
		}

		logger.Info("GOT", "tid", tid, "archetypes", archetypes)

		components.Dashboard().Render(r.Context(), w)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	r.NotFound(NotFoundHandler)

	s := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down")
		if err := s.Shutdown(context.Background()); err != nil {
			logger.Error("closing server",
				slog.Any("error", err),
			)
		}
	}()

	logger.Info(
		"starting server",
		slog.String("port", config.Port),
		slog.String("host", config.Host),
	)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	err := errors.New("Page Not Found")

	w.WriteHeader(http.StatusNotFound)
	switch {
	// case strings.Contains(accept, "text/html"):
	// 	w.Header().Set("Content-Type", "text/html")
	// 	page.Error(err).Render(r.Context(), w)
	case strings.Contains(accept, "text/plain"):
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))
	case strings.Contains(accept, "application/json"):
		fallthrough
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"error": err})
	}
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			logger.Info("replied to request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.statusCode,
				"duration", time.Since(start).String(),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
