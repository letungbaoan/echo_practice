package utils

import (
	"testing"
)

const testSecret = "test-secret-key"

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(42, testSecret)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected UserID=42, got %d", claims.UserID)
	}
	if claims.Issuer != "echo_practice" {
		t.Fatalf("unexpected issuer: %s", claims.Issuer)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, _ := GenerateToken(1, testSecret)
	if _, err := ParseToken(token, "another-secret"); err == nil {
		t.Fatal("expected error when parsing with wrong secret")
	}
}

func TestParseTokenGarbage(t *testing.T) {
	if _, err := ParseToken("not.a.jwt", testSecret); err == nil {
		t.Fatal("expected error parsing garbage")
	}
}
