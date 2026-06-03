package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
