package utils

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

func HandelError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.InvalidInput):
		http.Error(w, "bad request", http.StatusBadRequest)
	case errors.Is(err, entity.BadCredentials):
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	case errors.Is(err, entity.NotFoundError):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, entity.JWTError):
		http.Error(w, "invalid jwt", http.StatusUnauthorized)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
