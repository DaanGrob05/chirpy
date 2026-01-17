package handlers

import (
	"fmt"
	"net/http"

	apiconfig "example.com/chirpy/api_config"
)

func ServerMetricsHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func ResetHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Store(0)

		err := cfg.DbQueries.ResetUsers(r.Context())
		if err != nil {
			fmt.Println(err)
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
