package handlers

import (
	"fmt"
	"strings"
)

const MaxChirpLength = 140

func validateChirp(body string) (string, error) {

	if body == "" {
		return "", fmt.Errorf("Body cannot be empty")
	}

	if len(body) > MaxChirpLength {
		return "", fmt.Errorf("Chirp is too long")
	}

	return sanitiseInput(body), nil

}

func sanitiseInput(input string) string {

	words := strings.Split(input, " ")

	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

	for i, word := range words {
		lower := strings.ToLower(word)
		if _, ok := badWords[lower]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
