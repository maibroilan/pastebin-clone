package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maibroilan/pastebin-clone/server/internal/service"
)

func TestWriteJSON_Success(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{"message": "hello"}
	WriteJSON(w, 200, data)

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["message"] != "hello" {
		t.Errorf("expected message 'hello', got '%s'", result["message"])
	}
}

func TestWriteJSON_DifferentStatus(t *testing.T) {
	w := httptest.NewRecorder()

	WriteJSON(w, 201, map[string]string{"status": "created"})

	if w.Code != 201 {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestHandleError_PasteNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, service.ErrPasteNotFound)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleError_PasteExpired(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, service.ErrPasteExpired)

	if w.Code != http.StatusGone {
		t.Errorf("expected 410, got %d", w.Code)
	}
}

func TestHandleError_InvalidSlug(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, service.ErrInvalidSlug)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleError_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, service.ErrBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleError_NotAuthorized(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, service.ErrNotAuthorized)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleError_WrongPassword(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, service.ErrWrongPassword)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleError_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()

	handleError(w, errors.New("something weird"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
