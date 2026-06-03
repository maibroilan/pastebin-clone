package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSlug() (string, error) {
	b := make([]byte, 6)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GeneratePasswordHash(password *string) (*string, error) {
	if password == nil {
		return nil, nil
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(*password),
		bcrypt.DefaultCost,
	)

	hashString := string(hash)
	return &hashString, err
}

func CompareHashes(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	if err != nil {
		return false
	}

	return true
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
