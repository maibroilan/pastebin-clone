package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maibroilan/pastebin-clone/server/internal/db"
	"github.com/maibroilan/pastebin-clone/server/internal/handlers"
	"github.com/maibroilan/pastebin-clone/server/internal/service"
)

func main() {
	slog.Info("Starting Server...")
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://maibroilan:admin@localhost:5432/pastebin?sslmode=disable")
	if err != nil {
		slog.Error("couldn't initialize db pool", "error", err)
		os.Exit(1)
	}

	queries := db.New(pool)

	r := chi.NewRouter()

	// 🧰 global middleware stack
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.BodyLimit(1 << 20))
	r.Use(cors.Handler(cors.Options{
		// List of allowed origins. Use "*" for development only.
		AllowedOrigins: []string{"http://localhost:5173"}, // Your SvelteKit dev server
		// Allowed HTTP methods
		AllowedMethods: []string{"GET", "POST"},
		// Allowed headers (include Authorization for tokens, Content-Type for JSON)
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Paste-Password"},
		// Headers the browser is allowed to access
		ExposedHeaders: []string{"Link"},
		// Allow cookies to be sent/received
		AllowCredentials: false,
		// Cache preflight request result for 5 minutes
		MaxAge: 300,
	}))

	// 📦 API routes
	r.Route("/pastes", func(r chi.Router) {
		p := service.NewPasteService(*queries)
		h := handlers.NewPasteHandler(p)

		r.Post("/", h.CreatePaste)
		r.Get("/{slug}", h.GetPaste)
	})

	slog.Info("Server Started on localhost:8080")
	http.ListenAndServe(":8080", r)
}
