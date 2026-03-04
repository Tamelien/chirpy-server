package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	tokenSecret := "correct horse battery staple"
	testUserID := uuid.New()

	token, err := MakeJWT(testUserID, tokenSecret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	userID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}
	if userID != testUserID {
		t.Fatalf("expected password to match hash")
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatalf("expected error for wrong secret, got none")
	}

	token, err = MakeJWT(testUserID, tokenSecret, time.Nanosecond)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}
	time.Sleep(time.Millisecond)
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("expected error for expired tokens, got none")
	}
}
