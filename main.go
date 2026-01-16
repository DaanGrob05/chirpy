package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"example.com/chirpy/handlers"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(cfg.fileserverHits.Load())
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg apiConfig) getFileserverHits() int32 {
	return cfg.fileserverHits.Load()
}

func main() {
	mux := http.NewServeMux()
	apiCgf := apiConfig{}

	mux.Handle("GET /app", apiCgf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
	mux.Handle("GET /app/assets/", apiCgf.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets/")))))

	mux.HandleFunc("GET /api/healthz", handlers.GETHealthzHandler)

	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		hits := apiCgf.getFileserverHits()
		text := fmt.Sprintf("Hits: %v", hits)

		w.Write([]byte(text))
	})
	mux.HandleFunc("POST /api/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCgf.fileserverHits.Store(0)

		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
