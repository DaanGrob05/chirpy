package apiconfig

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"example.com/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetFileserverHits() int32 {
	return cfg.FileserverHits.Load()
}

func (cfg *ApiConfig) ServerMetricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.GetFileserverHits()
	text := fmt.Sprintf(`
			<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>
		`, hits)

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(text))
}

func (cfg *ApiConfig) ServerMetricsResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
}
