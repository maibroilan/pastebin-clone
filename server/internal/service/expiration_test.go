package service

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCalculateExpiration_ValidKeys(t *testing.T) {
	tests := []struct {
		key      string
		expected time.Duration
	}{
		{"1m", 1 * time.Minute},
		{"30m", 30 * time.Minute},
		{"1h", 1 * time.Hour},
		{"1d", 24 * time.Hour},
		{"7d", 7 * 24 * time.Hour},
		{"10d", 10 * 24 * time.Hour},
		{"30d", 30 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			before := time.Now()
			result, err := CalculateExpiration(tt.key)
			after := time.Now()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if !result.Valid {
				t.Fatal("expected Valid to be true")
			}

			expectedTime := before.Add(tt.expected)
			if result.Time.Before(expectedTime.Add(-time.Second)) || result.Time.After(after.Add(tt.expected).Add(time.Second)) {
				t.Errorf("expected time near %v, got %v", expectedTime, result.Time)
			}
		})
	}
}

func TestCalculateExpiration_InvalidKey(t *testing.T) {
	result, err := CalculateExpiration("invalid")

	if err != ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}

	if result != nil {
		t.Fatal("expected nil result for invalid key")
	}
}

func TestCalculateExpiration_EmptyString(t *testing.T) {
	result, err := CalculateExpiration("")

	if err != ErrBadRequest {
		t.Fatalf("expected ErrBadRequest for empty string, got %v", err)
	}

	if result != nil {
		t.Fatal("expected nil result for empty string")
	}
}

func TestIsExpired_Past(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	ts := PastTimestamp(past)

	if !IsExpired(ts) {
		t.Error("expected past timestamp to be expired")
	}
}

func TestIsExpired_Future(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	ts := FutureTimestamp(future)

	if IsExpired(ts) {
		t.Error("expected future timestamp to not be expired")
	}
}

func TestIsExpired_Now(t *testing.T) {
	now := time.Now()
	ts := PastTimestamp(now)

	if !IsExpired(ts) {
		t.Error("expected current time to be considered expired (After is strictly after)")
	}
}

func PastTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func FutureTimestamp(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
