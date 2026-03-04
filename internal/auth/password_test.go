package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	pw := "correct horse battery staple"

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}

	ok, err := CheckPasswordHash(pw, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash error: %v", err)
	}
	if !ok {
		t.Fatalf("expected password to match hash")
	}

	ok, err = CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash error: %v", err)
	}
	if ok {
		t.Fatalf("expected password NOT to match hash")
	}
}
