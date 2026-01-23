package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")

	bearer = strings.TrimPrefix(bearer, "Bearer ")

	if bearer == "" {
		return "", errors.New("Unauthorized.")
	}

	return bearer, nil
}
