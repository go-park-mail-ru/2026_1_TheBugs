package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetAccessToken(r *http.Request) (string, error) {
	const BearerPrefix = "Bearer "
	cred := r.Header.Get("Authorization")

	if cred == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(cred, BearerPrefix) {
		return "", errors.New("invalid authorization header: missing Bearer prefix")
	}
	token := strings.TrimPrefix(cred, BearerPrefix)
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}
