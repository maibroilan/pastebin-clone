package service

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestGeneratePasswordHash_Nil(t *testing.T) {
	hash, err := GeneratePasswordHash(nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hash != nil {
		t.Fatal("expected nil hash for nil password")
	}
}

func TestGeneratePasswordHash_NonNil(t *testing.T) {
	password := "mypassword"
	hash, err := GeneratePasswordHash(&password)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hash == nil {
		t.Fatal("expected non-nil hash")
	}

	if *hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if *hash == password {
		t.Fatal("hash should not equal plaintext password")
	}
}

func TestCompareHashes_CorrectPassword(t *testing.T) {
	password := "secret123"
	hash, err := GeneratePasswordHash(&password)
	if err != nil {
		t.Fatalf("unexpected error generating hash: %v", err)
	}

	if !CompareHashes(*hash, password) {
		t.Error("expected hashes to match for correct password")
	}
}

func TestCompareHashes_WrongPassword(t *testing.T) {
	password := "secret123"
	hash, err := GeneratePasswordHash(&password)
	if err != nil {
		t.Fatalf("unexpected error generating hash: %v", err)
	}

	if CompareHashes(*hash, "wrongpassword") {
		t.Error("expected hashes to not match for wrong password")
	}
}

func TestIsUniqueViolation_PgError23505(t *testing.T) {
	err := &pgconn.PgError{Code: "23505"}

	if !IsUniqueViolation(err) {
		t.Error("expected true for PgError with code 23505")
	}
}

func TestIsUniqueViolation_OtherPgErrorCode(t *testing.T) {
	err := &pgconn.PgError{Code: "23503"}

	if IsUniqueViolation(err) {
		t.Error("expected false for PgError with non-23505 code")
	}
}

func TestIsUniqueViolation_RegularError(t *testing.T) {
	err := errors.New("some error")

	if IsUniqueViolation(err) {
		t.Error("expected false for regular error")
	}
}

func TestIsUniqueViolation_Nil(t *testing.T) {
	if IsUniqueViolation(nil) {
		t.Error("expected false for nil error")
	}
}
