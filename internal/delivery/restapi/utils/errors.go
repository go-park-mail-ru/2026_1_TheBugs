package restapi

import (
	"encoding/json"
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
	switch err {
	case entity.InvalidInput:
		WriteError(w, "bad request", http.StatusBadRequest)
	case entity.BadCredentials:
		WriteError(w, "unauthorized", http.StatusUnauthorized)
	case entity.NotFoundError:
		WriteError(w, "not found", http.StatusNotFound)
	default:
		WriteError(w, "internal", http.StatusInternalServerError)
	}
}
