package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	apiconfig "example.com/chirpy/api_config"
	"example.com/chirpy/internal/auth"
	"github.com/google/uuid"
)

func RefreshJWTHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type returnBody struct {
			Token string `json:"token"`
		}

		bearer, err := auth.GetTokenFromHeader(r.Header, "Bearer ")
		if err != nil {
			returnError(err, w, http.StatusUnauthorized)
			return
		}

		userId, err := cfg.DbQueries.GetUserIdFromRefreshToken(r.Context(), bearer)
		if err != nil || userId == uuid.Nil {
			returnError(errors.New("Unauthorized."), w, http.StatusUnauthorized)
			return
		}

		jwt, err := auth.MakeJWT(userId, cfg.Secret)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		body := returnBody{
			Token: jwt,
		}

		jsonData, err := json.Marshal(body)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func RevokeJWTHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type returnBody struct {
			Token string `json:"token"`
		}

		bearer, err := auth.GetTokenFromHeader(r.Header, "Bearer ")
		if err != nil {
			returnError(err, w, http.StatusUnauthorized)
			return
		}

		err = cfg.DbQueries.RevokeRefreshToken(r.Context(), bearer)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
