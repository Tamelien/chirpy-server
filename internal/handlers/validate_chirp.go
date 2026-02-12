package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
		Valid bool `json:"valid"`
	}

	respBody := returnVals{
		Valid: true,
	}

	respondWithJSON(w, http.StatusOK, respBody)
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
