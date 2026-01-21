package apiconfig

import (
	"net/http"
	"sync/atomic"

	"example.com/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Secret         string
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
