package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tamelien/chirpy-server/internal/api"
)

func HandlerCreateUsers(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		params.Email = strings.TrimSpace(params.Email)
		if params.Email == "" {
			respondWithError(w, http.StatusBadRequest, "No email received.")
			return
		}

		user, err := cfg.DBQueries.CreateUser(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		type User struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}

		respBody := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		respondWithJSON(w, http.StatusCreated, respBody)

	}
}
