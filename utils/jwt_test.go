package utils

import "testing"

func TestGenerateAndVerifyToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateToken("user@example.com", 42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	userID, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}

	if userID != 42 {
		t.Fatalf("VerifyToken() userID = %d, want 42", userID)
	}
}

func TestVerifyTokenRejectsTokenWithDifferentSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-1")

	token, err := GenerateToken("user@example.com", 7)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	t.Setenv("JWT_SECRET", "secret-2")
	if _, err := VerifyToken(token); err == nil {
		t.Fatal("VerifyToken() expected error for mismatched secret")
	}
}
