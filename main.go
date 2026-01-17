package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	apiconfig "example.com/chirpy/api_config"
	"example.com/chirpy/handlers"
	"example.com/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can not connect to database.")
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	apiCgf := apiconfig.ApiConfig{
		DbQueries: dbQueries,
	}

	mux.Handle("GET /app", apiCgf.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
	mux.Handle("GET /app/assets/", apiCgf.MiddlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets/")))))

	mux.HandleFunc("GET /api/healthz", handlers.GETHealthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", handlers.POSTValidateChirp)

	mux.HandleFunc("GET /admin/metrics", handlers.ServerMetricsHandler(&apiCgf))
	mux.HandleFunc("POST /admin/reset", handlers.ResetHandler(&apiCgf))

	mux.HandleFunc("POST /api/users", handlers.CreateUserHandler(&apiCgf))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
