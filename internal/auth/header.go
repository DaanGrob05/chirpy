package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetTokenFromHeader(headers http.Header, prefix string) (string, error) {
	bearer := headers.Get("Authorization")

	if prefix[len(prefix)-1:] != " " {
		prefix = prefix + " "
	}

	bearer = strings.TrimPrefix(bearer, prefix)

	if bearer == "" {
		return "", errors.New("Unauthorized.")
	}

	return bearer, nil
}
