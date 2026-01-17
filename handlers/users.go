package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	apiconfig "example.com/chirpy/api_config"
)

func CreateUserHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
		}

		fmt.Println("Creating User")

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			returnError(err, w, http.StatusBadRequest)
			return
		}

		user, err := cfg.DbQueries.CreateUser(r.Context(), params.Email)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(user)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}
