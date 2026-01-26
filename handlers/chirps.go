package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	apiconfig "example.com/chirpy/api_config"
	"example.com/chirpy/internal/auth"
	"example.com/chirpy/internal/database"
	"example.com/chirpy/logging"
	"github.com/google/uuid"
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

		token, err := auth.GetTokenFromHeader(r.Header, "Bearer ")
		if err != nil {
			returnError(err, w, http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			returnError(errors.New("Unauthorized."), w, http.StatusUnauthorized)
			return
		}

		createParams.UserID = userId

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

		authorIdString := r.URL.Query().Get("author_id")
		sortOrder := r.URL.Query().Get("sort")

		var data []database.Chirp
		var err error

		if authorIdString == "" {
			data, err = cfg.DbQueries.GetChirps(r.Context())
			if err != nil {
				returnError(err, w, http.StatusInternalServerError)
				return
			}
		} else {
			authorId, err := uuid.Parse(authorIdString)
			if err != nil {
				returnError(err, w, http.StatusBadRequest)
				return
			}
			data, err = cfg.DbQueries.GetChirpsByUser(r.Context(), authorId)
			if err != nil {
				returnError(err, w, http.StatusInternalServerError)
				return
			}
		}

		slices.SortFunc(data, func(a, b database.Chirp) int {
			cmp := a.CreatedAt.Compare(b.CreatedAt)
			if sortOrder == "desc" {
				return -cmp
			}
			return cmp
		})

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

func GetOneChirpHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Log("Getting One Chirp")

		idString := r.PathValue("chirpID")
		if idString == "" {
			returnError(errors.New("Incorrect or empty id value."), w, http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idString)
		if err != nil {
			returnError(err, w, http.StatusNotFound)
			return
		}

		chirp, err := cfg.DbQueries.GetOneChirp(r.Context(), id)
		if err != nil {
			returnError(err, w, http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(chirp)
		if err != nil {
			returnError(err, w, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func DeleteChirpHandler(cfg *apiconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpIdString := r.PathValue("chirpID")
		if chirpIdString == "" {
			returnError(errors.New("Incorrect or empty id value."), w, http.StatusBadRequest)
			return
		}

		chirpId, err := uuid.Parse(chirpIdString)
		if err != nil {
			returnError(err, w, http.StatusNotFound)
			return
		}

		bearer, err := auth.GetTokenFromHeader(r.Header, "Bearer ")
		if err != nil {
			returnError(err, w, http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateJWT(bearer, cfg.Secret)
		if err != nil {
			returnError(err, w, http.StatusUnauthorized)
			return
		}

		chirp, err := cfg.DbQueries.GetOneChirp(r.Context(), chirpId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				returnError(errors.New("Chirp not found."), w, http.StatusNotFound)
				return
			}

			returnError(err, w, http.StatusInternalServerError)
			return
		}

		if chirp.UserID != userId {
			returnError(errors.New("Forbidden."), w, http.StatusForbidden)
			return
		}

		deleteParams := database.DeleteChirpParams{
			ID:     chirpId,
			UserID: userId,
		}
		err = cfg.DbQueries.DeleteChirp(r.Context(), deleteParams)
		if err != nil {
			returnError(err, w, http.StatusForbidden)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
