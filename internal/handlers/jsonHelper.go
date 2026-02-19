package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

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
