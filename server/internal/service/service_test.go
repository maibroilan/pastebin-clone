package service

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maibroilan/pastebin-clone/server/internal/db"
	"github.com/maibroilan/pastebin-clone/server/internal/model"
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

func TestPasteService_Create_Success(t *testing.T) {
	store := &FakePasteStore{
		CreatePasteFn: func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
			return db.Paste{
				Slug:      arg.Slug,
				Content:   arg.Content,
				ExpiresAt: arg.ExpiresAt,
			}, nil
		},
	}

	service := NewPasteService(store)

	password := "secret123"

	req := model.CreatePasteRequest{
		Content:    "test content",
		Password:   &password,
		Expiration: "1h",
	}

	paste, err := service.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if paste == nil || paste.Slug == "" || paste.Content != req.Content {
		t.Error("invalid paste returned")
	}
}

func TestPasteService_Create_RetryOnDuplicate(t *testing.T) {
	attempts := 0

	store := &FakePasteStore{
		CreatePasteFn: func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
			attempts++
			if attempts < 3 {
				// Simulate unique violation on first two tries
				return db.Paste{}, &pgconn.PgError{Code: "23505"}
			}
			return db.Paste{Slug: arg.Slug, Content: arg.Content}, nil
		},
	}

	service := NewPasteService(store)
	req := model.CreatePasteRequest{Content: "hello", Expiration: "1h"}

	paste, err := service.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if paste.Slug == "" {
		t.Error("paste should have slug")
	}
}

func TestPasteService_Create_FailsAfterMaxRetries(t *testing.T) {
	store := &FakePasteStore{
		CreatePasteFn: func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
			return db.Paste{}, &pgconn.PgError{Code: "23505"}
		},
	}

	service := NewPasteService(store)
	_, err := service.Create(context.Background(), model.CreatePasteRequest{Content: "test", Expiration: "1h"})
	if !errors.Is(err, ErrSlugGenFailed) {
		t.Errorf("expected ErrSlugGenFailed, got %v", err)
	}
}

func TestPasteService_Create_OtherError(t *testing.T) {
	expectedErr := errors.New("database down")

	store := &FakePasteStore{
		CreatePasteFn: func(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error) {
			return db.Paste{}, expectedErr
		},
	}

	service := NewPasteService(store)
	_, err := service.Create(context.Background(), model.CreatePasteRequest{Content: "test", Expiration: "1h"})
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected original error, got %v", err)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
func TestSlugGeneration(t *testing.T) {
	var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{6,20}$`)

	slug, err := GenerateSlug()

	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	if len(slug) != 8 {
		t.Fatalf("expected length 8 but got %d", len(slug))
	}

	if !slugRegex.MatchString(slug) {
		t.Fatalf("generated slug is not url safe")
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Get tests

func TestPasteService_Get_Success(t *testing.T) {
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:      slug,
				Content:   "hello world",
				ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			}, nil
		},
	}

	svc := NewPasteService(store)

	paste, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if paste.Content != "hello world" {
		t.Errorf("expected content 'hello world', got '%s'", paste.Content)
	}

	if paste.Slug != "abc123" {
		t.Errorf("expected slug 'abc123', got '%s'", paste.Slug)
	}
}

func TestPasteService_Get_EmptySlug(t *testing.T) {
	svc := NewPasteService(&FakePasteStore{})

	_, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: ""})
	if !errors.Is(err, ErrInvalidSlug) {
		t.Errorf("expected ErrInvalidSlug, got %v", err)
	}
}

func TestPasteService_Get_NotFound(t *testing.T) {
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{}, pgx.ErrNoRows
		},
	}

	svc := NewPasteService(store)

	_, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "nonexistent"})
	if !errors.Is(err, ErrPasteNotFound) {
		t.Errorf("expected ErrPasteNotFound, got %v", err)
	}
}

func TestPasteService_Get_DatabaseError(t *testing.T) {
	expectedErr := errors.New("connection lost")
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{}, expectedErr
		},
	}

	svc := NewPasteService(store)

	_, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "abc123"})
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected original error, got %v", err)
	}
}

func TestPasteService_Get_Expired(t *testing.T) {
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:      slug,
				Content:   "expired content",
				ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(-1 * time.Hour), Valid: true},
			}, nil
		},
	}

	svc := NewPasteService(store)

	_, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "abc123"})
	if !errors.Is(err, ErrPasteExpired) {
		t.Errorf("expected ErrPasteExpired, got %v", err)
	}
}

func TestPasteService_Get_PasswordRequired(t *testing.T) {
	hash := "$2a$10$somehash"
	store := &FakePasteStore{
		GetPasteFn: func(ctx context.Context, slug string) (db.Paste, error) {
			return db.Paste{
				Slug:         slug,
				Content:      "protected content",
				PasswordHash: &hash,
				ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			}, nil
		},
	}

	svc := NewPasteService(store)

	_, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "abc123"})
	if !errors.Is(err, ErrNotAuthorized) {
		t.Errorf("expected ErrNotAuthorized when no password provided, got %v", err)
	}
}

func TestPasteService_Get_WrongPassword(t *testing.T) {
	hash, _ := GeneratePasswordHash(stringPtr("correctpassword"))
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

	svc := NewPasteService(store)

	_, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "abc123", Password: "wrongpassword"})
	if !errors.Is(err, ErrWrongPassword) {
		t.Errorf("expected ErrWrongPassword, got %v", err)
	}
}

func TestPasteService_Get_CorrectPassword(t *testing.T) {
	hash, _ := GeneratePasswordHash(stringPtr("correctpassword"))
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

	svc := NewPasteService(store)

	paste, err := svc.Get(context.Background(), model.GetPasteRequest{Slug: "abc123", Password: "correctpassword"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if paste.Content != "protected content" {
		t.Errorf("expected content 'protected content', got '%s'", paste.Content)
	}
}

func stringPtr(s string) *string {
	return &s
}

// /////////////////////////////////////////////////////////////////////////////////
// PingDB tests

func TestPasteService_PingDB_Success(t *testing.T) {
	store := &FakePasteStore{
		PingFn: func(ctx context.Context) (int32, error) {
			return 1, nil
		},
	}

	svc := NewPasteService(store)

	res, err := svc.PingDB(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Status != "ready" {
		t.Errorf("expected status 'ready', got '%s'", res.Status)
	}

	if res.Checks["postgres"] != "up" {
		t.Errorf("expected postgres 'up', got '%s'", res.Checks["postgres"])
	}
}

func TestPasteService_PingDB_Failure(t *testing.T) {
	store := &FakePasteStore{
		PingFn: func(ctx context.Context) (int32, error) {
			return 0, errors.New("connection refused")
		},
	}

	svc := NewPasteService(store)

	res, err := svc.PingDB(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if res.Status != "not_ready" {
		t.Errorf("expected status 'not_ready', got '%s'", res.Status)
	}

	if res.Checks["postgres"] != "down" {
		t.Errorf("expected postgres 'down', got '%s'", res.Checks["postgres"])
	}
}
