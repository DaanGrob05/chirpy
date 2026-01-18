package handlers

import (
	"encoding/json"
	"net/http"

	apiconfig "example.com/chirpy/api_config"
	"example.com/chirpy/internal/database"
	"example.com/chirpy/logging"
)

func CreateChirpHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Log("Creating Chirp")

		decoder := json.NewDecoder(r.Body)
		createParams := database.CreateChirpParams{}
		err := decoder.Decode(&createParams)
		if err != nil {
			returnError(err, w, http.StatusBadRequest)
			return
		}

		createdChirp, err := cfg.DbQueries.CreateChirp(r.Context(), createParams)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(createdChirp)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	}
}

func GetChirpsHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Log("Getting Chirps")

		data, err := cfg.DbQueries.GetChirps(r.Context())
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
