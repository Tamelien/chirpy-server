package main

import (
	"log"
	"net/http"

	"github.com/tamelien/chirpy-server/internal/api"
	"github.com/tamelien/chirpy-server/internal/handlers"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	cfg := &api.ApiConfig{}
	fileServer := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(filepathRoot)),
	)

	mux.Handle("/app/", handlers.HandlerMetricsInc(cfg)(fileServer))
	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("GET  /admin/metrics", handlers.HandlerMetricsRead(cfg))
	mux.HandleFunc("POST /admin/reset", handlers.HandlerReset(cfg))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
