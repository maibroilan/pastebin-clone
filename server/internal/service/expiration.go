package service

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var ExpirationMappings = map[string]time.Duration{
	"1m":  1 * time.Minute,
	"30m": 30 * time.Minute,
	"1h":  1 * time.Hour,
	"1d":  24 * time.Hour,
	"7d":  7 * 24 * time.Hour,
	"10d": 10 * 24 * time.Hour,
	"30d": 30 * 24 * time.Hour,
}

func CalculateExpiration(expiration string) (*pgtype.Timestamptz, error) {
	duration, ok := ExpirationMappings[expiration]
	if !ok {
		return nil, ErrBadRequest
	}

	return &pgtype.Timestamptz{
		Time:  time.Now().Add(duration),
		Valid: true,
	}, nil
}

func IsExpired(expiration pgtype.Timestamptz) bool {
	return time.Now().After(expiration.Time)
}
