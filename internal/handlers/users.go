package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tamelien/chirpy-server/internal/api"
	"github.com/tamelien/chirpy-server/internal/auth"
	"github.com/tamelien/chirpy-server/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // ignored by json.Marshal
}

type parametersLogin struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func HandlerCreateUsers(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		params := parametersLogin{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		params.Email = strings.TrimSpace(params.Email)
		if params.Email == "" {
			respondWithError(w, http.StatusBadRequest, "No Email received.")
			return
		}

		params.Password = strings.TrimSpace(params.Password)
		if params.Password == "" {
			respondWithError(w, http.StatusBadRequest, "No Password received.")
			return
		}

		hash_password, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := cfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hash_password})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
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

func HandlerLogin(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		decoder := json.NewDecoder(req.Body)
		params := parametersLogin{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding login user: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		params.Email = strings.TrimSpace(params.Email)
		if params.Email == "" {
			respondWithError(w, http.StatusBadRequest, "No email received.")
			return
		}
		user, err := cfg.DBQueries.GetUser(req.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil || ok == false {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		respBody := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		respondWithJSON(w, http.StatusOK, respBody)

	}
}
