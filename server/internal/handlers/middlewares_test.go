package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBodyLimit_LimitApplied(t *testing.T) {
	var receivedBody io.ReadCloser

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody = r.Body
		w.WriteHeader(200)
	})

	handler := BodyLimit(100)(inner)

	body := strings.Repeat("a", 50)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if receivedBody == nil {
		t.Fatal("expected body to be set")
	}

	// Read all — should succeed since body is under limit
	data, err := io.ReadAll(receivedBody)
	if err != nil {
		t.Fatalf("unexpected error reading body: %v", err)
	}

	if len(data) != 50 {
		t.Errorf("expected 50 bytes, got %d", len(data))
	}
}

func TestBodyLimit_OverLimit(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read more than the limit
		buf := make([]byte, 200)
		_, err := io.ReadFull(r.Body, buf)
		if err != nil {
			http.Error(w, "body too large", http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(200)
	})

	handler := BodyLimit(100)(inner)

	body := strings.Repeat("a", 150)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// MaxBytesReader causes an error when reading past the limit
	if w.Code == http.StatusOK {
		t.Error("expected non-200 status when body exceeds limit")
	}
}

func TestBodyLimit_ChainOfHandlers(t *testing.T) {
	var nextCalled bool

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(200)
	})

	handler := BodyLimit(1024)(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("expected next handler to be called")
	}

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
