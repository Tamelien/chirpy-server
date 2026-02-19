package handlers

import (
	"log"
	"net/http"

	"github.com/tamelien/chirpy-server/internal/api"
)

func HandlerReset(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if cfg.PLATFORM != "dev" {
			http.Error(w, "No Permission", http.StatusForbidden)
			return
		}

		err := cfg.DBQueries.Reset(req.Context())
		if err != nil {
			log.Printf("Error reseting Database: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		cfg.FileserverHits.Store(0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}
}
