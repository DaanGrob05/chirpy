package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
)

func POSTValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validBody struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		returnError(err, w, http.StatusInternalServerError)
	}

	if len(params.Body) > 140 {
		returnError(errors.New("Chirp is too long."), w, http.StatusBadRequest)
		return
	}

	sanitizedText := checkForProfanity(params.Body)

	returnBody := validBody{
		Valid:       true,
		CleanedBody: sanitizedText,
	}

	data, err := json.Marshal(returnBody)
	if err != nil {
		returnError(err, w, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func returnError(err error, w http.ResponseWriter, httpCode int) {
	type errorBody struct {
		Error string `json:"error"`
	}

	body := errorBody{
		Error: "Something went wrong. Please try again.",
	}

	data, marshallErr := json.Marshal(body)
	if marshallErr != nil {
		returnError(marshallErr, w, http.StatusInternalServerError)
		return
	}

	log.Printf("Error occurred while decoding parameters: %v", err)
	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return
}

func checkForProfanity(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")

	for index, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[index] = "****"
		}
	}

	return strings.Join(words, " ")
}
