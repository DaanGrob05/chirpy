package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", errors.New("Unauthorized.")
	}

	bearer = strings.TrimPrefix(bearer, "Bearer ")

	return bearer, nil
}
