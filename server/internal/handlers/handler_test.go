package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maibroilan/pastebin-clone/server/internal/db"
	"github.com/maibroilan/pastebin-clone/server/internal/model"
	"github.com/maibroilan/pastebin-clone/server/internal/service"
)

type FakePasteStore struct {
	CreatePasteFn func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error)
	GetPasteFn    func(ctx context.Context, slug string) (db.Paste, error)
	PingFn        func(ctx context.Context) (int32, error)
}

func (f *FakePasteStore) CreatePaste(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
	if f.CreatePasteFn != nil {
		return f.CreatePasteFn(ctx, arg)
	}
	return db.Paste{}, errors.New("CreatePasteFn not set")
}

func (f *FakePasteStore) GetPaste(ctx context.Context, slug string) (db.Paste, error) {
	if f.GetPasteFn != nil {
		return f.GetPasteFn(ctx, slug)
	}
	return db.Paste{}, errors.New("GetPasteFn not set")
}

func (f *FakePasteStore) Ping(ctx context.Context) (int32, error) {
	if f.PingFn != nil {
		return f.PingFn(ctx)
	}
	return 0, errors.New("PingFn not set")
}

func newTestHandler(store *FakePasteStore) *PasteHandler {
	svc := service.NewPasteService(store)
	return NewPasteHandler(svc)
}

// /////////////////////////////////////////////////////////////////////////////////
// CheckLive

func TestCheckLive(t *testing.T) {
	handler := newTestHandler(&FakePasteStore{})

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()

	handler.CheckLive(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp model.LiveCheckResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// CheckReady

func TestCheckReady_Success(t *testing.T) {
	store := &FakePasteStore{
		PingFn: func(ctx context.Context) (int32, error) {
			return 1, nil
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.CheckReady(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp model.ReadyCheckResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ready" {
		t.Errorf("expected status 'ready', got '%s'", resp.Status)
	}
}

func TestCheckReady_Down(t *testing.T) {
	store := &FakePasteStore{
		PingFn: func(ctx context.Context) (int32, error) {
			return 0, errors.New("db down")
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.CheckReady(w, req)

	if w.Code != 503 {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// CreatePaste

func TestCreatePaste_Success(t *testing.T) {
	store := &FakePasteStore{
		CreatePasteFn: func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
			return db.Paste{
				Slug:      arg.Slug,
				Content:   arg.Content,
				ExpiresAt: arg.ExpiresAt,
			}, nil
		},
	}
	handler := newTestHandler(store)

	body, _ := json.Marshal(model.CreatePasteRequest{
		Content:    "test content",
		Expiration: "1h",
	})

	req := httptest.NewRequest(http.MethodPost, "/paste", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreatePaste(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp model.CreatePasteResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Slug == "" {
		t.Error("expected non-empty slug")
	}
}

func TestCreatePaste_EmptyContent(t *testing.T) {
	handler := newTestHandler(&FakePasteStore{})

	body, _ := json.Marshal(model.CreatePasteRequest{
		Content:    "",
		Expiration: "1h",
	})

	req := httptest.NewRequest(http.MethodPost, "/paste", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreatePaste(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreatePaste_InvalidJSON(t *testing.T) {
	handler := newTestHandler(&FakePasteStore{})

	req := httptest.NewRequest(http.MethodPost, "/paste", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreatePaste(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreatePaste_DBError(t *testing.T) {
	store := &FakePasteStore{
		CreatePasteFn: func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
			return db.Paste{}, errors.New("database down")
		},
	}
	handler := newTestHandler(store)

	body, _ := json.Marshal(model.CreatePasteRequest{
		Content:    "test content",
		Expiration: "1h",
	})

	req := httptest.NewRequest(http.MethodPost, "/paste", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreatePaste(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// GetPaste

func TestGetPaste_Success(t *testing.T) {
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:      slug,
				Content:   "hello world",
				ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			}, nil
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/paste/abc123", nil)

	// Set chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler.GetPaste(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp model.GetPasteResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Content != "hello world" {
		t.Errorf("expected content 'hello world', got '%s'", resp.Content)
	}
}

func TestGetPaste_NotFound(t *testing.T) {
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{}, pgx.ErrNoRows
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/paste/nonexistent", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler.GetPaste(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetPaste_Expired(t *testing.T) {
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:      slug,
				Content:   "expired",
				ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(-1 * time.Hour), Valid: true},
			}, nil
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/paste/abc123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler.GetPaste(w, req)

	if w.Code != http.StatusGone {
		t.Errorf("expected 410, got %d", w.Code)
	}
}

func TestGetPaste_PasswordRequired(t *testing.T) {
	hash, _ := service.GeneratePasswordHash(ptrStr("secret"))
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:         slug,
				Content:      "protected",
				PasswordHash: hash,
				ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			}, nil
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/paste/abc123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler.GetPaste(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetPaste_WrongPassword(t *testing.T) {
	hash, _ := service.GeneratePasswordHash(ptrStr("secret"))
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:         slug,
				Content:      "protected",
				PasswordHash: hash,
				ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			}, nil
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/paste/abc123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	req.Header.Set("X-Paste-Password", "wrongpassword")

	w := httptest.NewRecorder()

	handler.GetPaste(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetPaste_CorrectPassword(t *testing.T) {
	hash, _ := service.GeneratePasswordHash(ptrStr("secret"))
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:         slug,
				Content:      "protected content",
				PasswordHash: hash,
				ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			}, nil
		},
	}
	handler := newTestHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/paste/abc123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	req.Header.Set("X-Paste-Password", "secret")

	w := httptest.NewRecorder()

	handler.GetPaste(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp model.GetPasteResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Content != "protected content" {
		t.Errorf("expected content 'protected content', got '%s'", resp.Content)
	}
}

func ptrStr(s string) *string {
	return &s
}
