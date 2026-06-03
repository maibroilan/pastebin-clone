package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Paste struct {
	Slug      string             `json:"slug"`
	Content   string             `json:"content"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

type CreatePasteRequest struct {
	Content    string  `json:"content"`
	Expiration string  `json:"expiration"` // "1h", "1d", etc.
	Password   *string `json:"password"`
}

type CreatePasteResponse struct {
	Slug      string    `json:"slug"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateGetResponse struct {
	Content   string    `json:"content"`
	ExpiresAt time.Time `json:"expires_at"`
}
