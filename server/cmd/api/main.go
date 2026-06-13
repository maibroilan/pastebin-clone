package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/maibroilan/pastebin-clone/server/internal/db"
	"github.com/maibroilan/pastebin-clone/server/internal/handlers"
	"github.com/maibroilan/pastebin-clone/server/internal/service"
)

const (
	defaultPort            = "8080"
	defaultDBTimeout       = 10 * time.Second
	defaultShutdownTimeout = 30 * time.Second
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		slog.Error("DB_URL environment variable is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultDBTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("couldn't initialize db pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("couldn't connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("database connection established")

	queries := db.New(pool)

	// Router

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.BodyLimit(1 << 20)) // 1 MiB

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173" // Default for local development
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Paste-Password"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	p := service.NewPasteService(queries)
	h := handlers.NewPasteHandler(p)

	r.Route("/pastes", func(r chi.Router) {
		r.Post("/", h.CreatePaste)
		r.Get("/{slug}", h.GetPaste)
	})

	r.Route("/livez", func(r chi.Router) {
		r.Get("/", h.CheckLive)
	})

	r.Route("/readyz", func(r chi.Router) {
		r.Get("/", h.CheckReady)
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("server failed to start", "error", err)
		os.Exit(1)

	case sig := <-shutdown:
		slog.Info("shutdown signal received", "signal", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			if err := server.Close(); err != nil {
				slog.Error("server close failed", "error", err)
			}
		} else {
			slog.Info("server gracefully stopped")
		}
	}
}
