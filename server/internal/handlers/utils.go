package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/maibroilan/pastebin-clone/server/internal/service"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("couldn't write json to output", "error", err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	switch {
	case errors.Is(err, service.ErrPasteNotFound):
		http.Error(w, "not found", http.StatusNotFound)

	case errors.Is(err, service.ErrPasteExpired):
		http.Error(w, "paste expired", http.StatusGone)

	case errors.Is(err, service.ErrInvalidSlug):
		http.Error(w, "invalid slug", http.StatusBadRequest)

	case errors.Is(err, service.ErrBadRequest):
		http.Error(w, "bad request", http.StatusBadRequest)

	case errors.Is(err, service.ErrNotAuthorized):
		http.Error(w, "not authorized", http.StatusUnauthorized)

	case errors.Is(err, service.ErrWrongPassword):
		http.Error(w, "wrong password", http.StatusUnauthorized)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
