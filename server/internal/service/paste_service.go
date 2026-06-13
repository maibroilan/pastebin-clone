package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/maibroilan/pastebin-clone/server/internal/db"
	"github.com/maibroilan/pastebin-clone/server/internal/model"
)

var (
	ErrPasteNotFound = errors.New("paste not found")
	ErrPasteExpired  = errors.New("paste expired")
	ErrInvalidSlug   = errors.New("invalid slug")
	ErrSlugGenFailed = errors.New("couldn't find a unique slug after 5 tries")
	ErrBadRequest    = errors.New("bad request")
	ErrNotAuthorized = errors.New("not authorized")
	ErrWrongPassword = errors.New("wrong password")
)

type PasteStore interface {
	CreatePaste(ctx context.Context, arg db.CreatePasteParams) (db.Paste, error)
	GetPaste(ctx context.Context, slug string) (db.Paste, error)
	Ping(ctx context.Context) (int32, error)
}

type PasteService struct {
	Queries PasteStore
}

func NewPasteService(ps PasteStore) *PasteService {
	return &PasteService{
		Queries: ps,
	}
}

func (s *PasteService) Create(ctx context.Context, paste model.CreatePasteRequest) (*model.Paste, error) {
	for range 5 {
		slug, err := GenerateSlug()
		if err != nil {
			slog.Error("couldn't create slug", "error", err)
			return nil, err
		}

		password_hash, err := GeneratePasswordHash(paste.Password)
		if err != nil {
			slog.Error("couldn't generate hash", "error", err)
			return nil, err
		}

		expirationDate, err := CalculateExpiration(paste.Expiration)
		if err != nil {
			slog.Error("couldn't calculate expiration date", "error", err)
			return nil, err
		}

		paste, err := s.Queries.CreatePaste(ctx, db.CreatePasteParams{
			Slug:         slug,
			Content:      paste.Content,
			PasswordHash: password_hash,
			ExpiresAt:    *expirationDate,
		})

		if err == nil {
			return &model.Paste{
				Slug:      paste.Slug,
				Content:   paste.Content,
				ExpiresAt: paste.ExpiresAt,
			}, nil
		}

		if !IsUniqueViolation(err) {
			slog.Error("couldn't create paste", "error", err)
			return nil, err
		}

		// try again
	}

	return nil, ErrSlugGenFailed
}

func (s *PasteService) Get(ctx context.Context, req model.GetPasteRequest) (*model.Paste, error) {
	if req.Slug == "" {
		return nil, ErrInvalidSlug
	}

	paste, err := s.Queries.GetPaste(ctx, req.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPasteNotFound
		}
		return nil, err
	}

	if IsExpired(paste.ExpiresAt) {
		return nil, ErrPasteExpired
	}

	if paste.PasswordHash != nil {
		if req.Password == "" {
			return nil, ErrNotAuthorized
		}

		if !CompareHashes(*paste.PasswordHash, req.Password) {
			return nil, ErrWrongPassword
		}
	}

	return &model.Paste{
		Slug:      req.Slug,
		Content:   paste.Content,
		ExpiresAt: paste.ExpiresAt,
	}, nil
}

func (s *PasteService) PingDB(ctx context.Context) (*model.ReadyCheckResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, PostgresErr := s.Queries.Ping(ctx)

	if PostgresErr != nil {
		return &model.ReadyCheckResponse{
			Status: "not_ready",
			Checks: map[string]string{
				"postgres": "down",
			},
		}, PostgresErr
	}

	return &model.ReadyCheckResponse{
		Status: "ready",
		Checks: map[string]string{
			"postgres": "up",
		},
	}, nil
}
