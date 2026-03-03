package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

func WriteError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.ErrorResponse{
		Error: msg,
	})
}

func HandelError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.InvalidInput):
		WriteError(w, "bad request", http.StatusBadRequest)
	case errors.Is(err, entity.BadCredentials):
		WriteError(w, "unauthorized", http.StatusUnauthorized)
	case errors.Is(err, entity.NotFoundError):
		WriteError(w, "not found", http.StatusNotFound)
	case errors.Is(err, entity.AlredyExitError):
		WriteError(w, "record alredy existed", http.StatusConflict)
	case errors.Is(err, entity.JWTError):
		WriteError(w, "bad jwt", http.StatusUnauthorized)
	case errors.Is(err, entity.OffsetOutOfRange):
		WriteError(w, "not enought records", http.StatusNotFound)
	default:
		WriteError(w, "internal", http.StatusInternalServerError)
	}
}
