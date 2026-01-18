package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	apiconfig "example.com/chirpy/api_config"
	"example.com/chirpy/internal/database"
)

func CreateChirpHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Creating Chirp")

		decoder := json.NewDecoder(r.Body)
		createParams := database.CreateChirpParams{}
		err := decoder.Decode(&createParams)
		if err != nil {
			returnError(err, w, http.StatusBadRequest)
		}

		fmt.Println(createParams)

		createdChirp, err := cfg.DbQueries.CreateChirp(r.Context(), createParams)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
		}

		data, err := json.Marshal(createdChirp)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	}
}
