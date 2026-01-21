package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/chirpy/logging"
)

func returnError(err error, w http.ResponseWriter, httpCode int) {
	logging.Log(fmt.Sprintf("Error: %v", err.Error()))
	type errorBody struct {
		Error string `json:"error"`
	}

	body := errorBody{
		Error: err.Error(),
	}

	data, marshallErr := json.Marshal(body)
	if marshallErr != nil {
		returnError(marshallErr, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return
}
