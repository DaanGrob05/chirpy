package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	apiconfig "example.com/chirpy/api_config"
	"example.com/chirpy/internal/auth"
	"example.com/chirpy/internal/database"
	"example.com/chirpy/logging"
	"github.com/google/uuid"
)

func CreateUserHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Log("Creating User")

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			returnError(err, w, http.StatusBadRequest)
			return
		}

		if params.Email == "" || params.Password == "" {
			returnError(errors.New("Email and password are required."), w, http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		createParams := database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		}

		user, err := cfg.DbQueries.CreateUser(r.Context(), createParams)
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

func LoginHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Log("Logging In")

		type parameters struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds int    `json:"expires_in_seconds"`
		}

		type exportBody struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
			Token     string    `json:"token"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			returnError(errors.New("Unauthorized."), w, http.StatusUnauthorized)
			return
		}

		user, err := cfg.DbQueries.GetUser(r.Context(), params.Email)
		if err != nil {
			returnError(errors.New("Unauthorized."), w, http.StatusUnauthorized)
			return
		}

		isValid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil || !isValid {
			returnError(errors.New("Unauthorized."), w, http.StatusUnauthorized)
			return
		}

		if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds >= 3600 {
			params.ExpiresInSeconds = 3600
		}

		duration, err := time.ParseDuration(fmt.Sprintf("%vs", params.ExpiresInSeconds))
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		token, err := auth.MakeJWT(user.ID, cfg.Secret, duration)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		body := exportBody{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		}

		jsonData, err := json.Marshal(body)
		if err != nil {
			returnError(err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
