package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tamelien/chirpy-server/internal/api"
	"github.com/tamelien/chirpy-server/internal/auth"
)

func HandlerPolkaWebhook(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		key, err := auth.GetAPIKey(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if key != cfg.POLKA_KEY {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type parameters struct {
			Event *string `json:"event"`
			Data  *struct {
				UserID *string `json:"user_id"`
			} `json:"data"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if params.Event == nil || *params.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if params.Data == nil || params.Data.UserID == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		id, err := uuid.Parse(*params.Data.UserID)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		_, err = cfg.DBQueries.UpdateUserIsChirpyRed(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
