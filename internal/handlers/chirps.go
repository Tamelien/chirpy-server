package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tamelien/chirpy-server/internal/api"
	"github.com/tamelien/chirpy-server/internal/database"
)

func HandlerCreateChirps(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		params.Body, err = validateChirp(params.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		chirp, err := cfg.DBQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body, UserID: params.UserID})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respBody := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		respondWithJSON(w, http.StatusCreated, respBody)

	}
}
