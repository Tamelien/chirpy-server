package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tamelien/chirpy-server/internal/api"
)

func HandlerGetChirp(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request")
			return
		}

		chirp, err := cfg.DBQueries.GetChirp(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "Chirp not found")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		respBody := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		respondWithJSON(w, http.StatusOK, respBody)

	}
}
