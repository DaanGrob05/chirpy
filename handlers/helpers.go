package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

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
