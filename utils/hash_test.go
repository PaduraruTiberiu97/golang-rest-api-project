package utils

import "testing"

func TestHashPasswordAndCheck(t *testing.T) {
	password := "strong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == password {
		t.Fatal("expected hashed password to differ from input")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("expected password/hash validation to succeed")
	}
}

func TestCheckPasswordHashRejectsWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if CheckPasswordHash("incorrect", hash) {
		t.Fatal("expected password/hash validation to fail")
	}
}
