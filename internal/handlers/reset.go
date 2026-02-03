package handlers

import (
	"net/http"

	"github.com/tamelien/chirpy-server/internal/api"
)

func HandlerReset(cfg *api.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset Fileserverhits count"))
		cfg.FileserverHits.Store(0)
	}
}
