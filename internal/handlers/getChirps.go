package handlers

import (
	"net/http"

	"github.com/tamelien/chirpy-server/internal/api"
)

func HandlerGetChirps(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cfg.DBQueries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		resp := make([]Chirp, 0, len(chirps))

		for _, chirp := range chirps {
			respBody := Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
			resp = append(resp, respBody)
		}

		respondWithJSON(w, http.StatusOK, resp)

	}
}
