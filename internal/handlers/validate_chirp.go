package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const MaxChirpLength = 140

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if params.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Body cannot be empty")
		return
	}

	if len(params.Body) > MaxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := returnVals{
		CleanedBody: sanitiseInput(params.Body),
	}

	respondWithJSON(w, http.StatusOK, respBody)
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

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respBody := returnVals{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
