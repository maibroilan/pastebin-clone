package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/maibroilan/pastebin-clone/server/internal/model"
	"github.com/maibroilan/pastebin-clone/server/internal/service"
)

type PasteHandler struct {
	service *service.PasteService
}

func NewPasteHandler(s *service.PasteService) *PasteHandler {
	return &PasteHandler{
		service: s,
	}
}

func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		1<<20, // 1 MiB
	)

	var req model.CreatePasteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("couldnt unmarshal request", "error", err)
		http.Error(w, "invalid request", 400)
		return
	}

	if req.Content == "" {
		handleError(w, service.ErrBadRequest)
		return
	}

	paste, err := h.service.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	WriteJSON(w, 200, model.CreatePasteResponse{
		Slug:      paste.Slug,
		ExpiresAt: paste.ExpiresAt.Time,
	})
}

func (h *PasteHandler) GetPaste(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	password := r.Header.Get("X-Paste-Password")

	paste, err := h.service.Get(
		r.Context(),
		model.GetPasteRequest{
			Slug:     slug,
			Password: password,
		},
	)

	if err != nil {
		handleError(w, err)
		return
	}

	WriteJSON(w, 200, model.GetPasteResponse{
		Content:   paste.Content,
		ExpiresAt: paste.ExpiresAt.Time,
	})
}
