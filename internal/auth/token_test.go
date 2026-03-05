package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		headerValue string
		expected    string
		expectError error
	}{
		{
			headerValue: "Bearer TOKEN_STRING",
			expected:    "TOKEN_STRING",
			expectError: nil,
		},
		{
			headerValue: "",
			expectError: errors.New("No Authorization Header"),
		},
		{
			headerValue: "Bearer",
			expectError: errors.New("Invalid Authorization Header"),
		},
		{
			headerValue: "Basic TOKEN_STRING",
			expectError: errors.New("No Bearer in Authorization Header"),
		},
	}

	for _, test := range tests {
		headers := http.Header{}
		headers.Set("Authorization", test.headerValue)

		token, err := GetBearerToken(headers)
		if err != nil {
			if err.Error() != test.expectError.Error() {
				t.Fatalf("expected %s, got %s", test.expectError, err)
			}
			continue
		}

		if token != test.expected {
			t.Fatalf("expected %s, got %s", test.expected, token)
		}
	}

}
