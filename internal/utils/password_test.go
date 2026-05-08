package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	plain := "s3cret-pass"

	hash, err := HashPassword(plain)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if hash == plain {
		t.Fatal("hash must differ from plain")
	}

	if !CheckPassword(plain, hash) {
		t.Fatal("expected check to succeed for correct password")
	}
	if CheckPassword("wrong-pass", hash) {
		t.Fatal("expected check to fail for wrong password")
	}
}

func TestHashIsRandomized(t *testing.T) {
	plain := "same-password"
	h1, _ := HashPassword(plain)
	h2, _ := HashPassword(plain)
	if h1 == h2 {
		t.Fatal("bcrypt should produce different hashes per call (random salt)")
	}
}
