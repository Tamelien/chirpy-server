package handlers

import (
	"database/sql"
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
	JWTToken       string    `json:"token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
	HashedPassword string    `json:"-"` // ignored by json.Marshal
}

type parametersLogin struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Expires  *int   `json:"expires_in_seconds"`
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
			log.Printf("Error Password hash: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := cfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hash_password})
		if err != nil {
			log.Printf("Error DB Query: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respBody := User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		}

		respondWithJSON(w, http.StatusCreated, respBody)

	}
}

func HandlerUpdateUser(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.SECRET)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		decoder := json.NewDecoder(req.Body)
		params := parametersLogin{}
		err = decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		params.Email = strings.TrimSpace(params.Email)
		if params.Email == "" {
			respondWithError(w, http.StatusUnauthorized, "No Email received.")
			return
		}

		params.Password = strings.TrimSpace(params.Password)
		if params.Password == "" {
			respondWithError(w, http.StatusUnauthorized, "No Password received.")
			return
		}

		hash_password, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("Error Password hash: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		user, err := cfg.DBQueries.UpdateUser(req.Context(), database.UpdateUserParams{ID: userID, Email: params.Email, HashedPassword: hash_password})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		respBody := User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		}

		respondWithJSON(w, http.StatusOK, respBody)
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

		expiresIn := time.Second * 3600 // Default expiration time of 1 hour.
		if params.Expires != nil && *params.Expires > 0 && *params.Expires < 3600 {
			expiresIn = time.Second * time.Duration(*params.Expires)
		}

		user, err := cfg.DBQueries.GetUser(req.Context(), params.Email)
		if err != nil {
			log.Printf("Error DB Query: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil {
			log.Printf("Error Check Passwordhash: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		if ok == false {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		token, err := auth.MakeJWT(user.ID, cfg.SECRET, expiresIn)
		if err != nil {
			log.Printf("Error MakeJWT: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error generating token")
			return
		}

		rtoken, err := cfg.DBQueries.CreateRefreshToken(req.Context(),
			database.CreateRefreshTokenParams{
				Token:     auth.MakeRefreshToken(),
				UserID:    user.ID,
				ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
				RevokedAt: sql.NullTime{},
			})
		if err != nil {
			log.Printf("Error Create Refresh Token: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error Create Refresh Token")
			return
		}

		respBody := User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			JWTToken:     token,
			RefreshToken: rtoken.Token,
			IsChirpyRed:  user.IsChirpyRed.Bool,
		}

		respondWithJSON(w, http.StatusOK, respBody)

	}
}

func HandlerRefresh(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rtoken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		user, err := cfg.DBQueries.GetUserFromRefreshToken(req.Context(), rtoken)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		token, err := auth.MakeJWT(user.ID, cfg.SECRET, time.Hour)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error generating token")
			return
		}

		respBody := struct {
			Token string `json:"token"`
		}{

			Token: token,
		}

		respondWithJSON(w, http.StatusOK, respBody)
	}
}

func HandlerRevoke(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rtoken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}

		err = cfg.DBQueries.RevokeRToken(req.Context(), rtoken)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Token")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
