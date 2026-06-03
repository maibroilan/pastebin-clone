package service

import (
	"context"
	"errors"
	"log/slog"

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
)

type PasteService struct {
	Queries db.Queries
}

func NewPasteService(q db.Queries) *PasteService {
	return &PasteService{
		Queries: q,
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

func (s *PasteService) Get(ctx context.Context, slug string) (*model.Paste, error) {
	if slug == "" {
		return nil, ErrInvalidSlug
	}

	paste, err := s.Queries.GetPaste(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPasteNotFound
		}
		return nil, err
	}

	if IsExpired(paste.ExpiresAt) {
		return nil, ErrPasteExpired
	}

	return &model.Paste{
		Slug:      slug,
		Content:   paste.Content,
		ExpiresAt: paste.ExpiresAt,
	}, nil
}
