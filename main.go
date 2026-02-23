package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tamelien/chirpy-server/internal/api"
	"github.com/tamelien/chirpy-server/internal/database"
	"github.com/tamelien/chirpy-server/internal/handlers"
)

func main() {
	cfg := &api.ApiConfig{}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM not set")
	}
	cfg.PLATFORM = platform

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	cfg.DBQueries = dbQueries

	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()

	fileServer := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(filepathRoot)),
	)

	mux.Handle("/app/", handlers.HandlerMetricsInc(cfg)(fileServer))

	//admin
	mux.HandleFunc("GET /admin/metrics", handlers.HandlerMetricsRead(cfg))
	mux.HandleFunc("POST /admin/reset", handlers.HandlerReset(cfg))

	//api
	mux.HandleFunc("GET /api/healthz", handlers.HealthHandler)
	mux.HandleFunc("GET /api/chirps", handlers.HandlerGetChirps(cfg))
	mux.HandleFunc("POST /api/chirps", handlers.HandlerCreateChirps(cfg))
	mux.HandleFunc("POST /api/users", handlers.HandlerCreateUsers(cfg))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
