package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

func HandelError(w http.ResponseWriter, err error) {
	switch err {
	case entity.InvalidInput:
		http.Error(w, "bad request data", http.StatusBadRequest)
	case entity.BadCredentials:
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	case entity.NotFoundError:
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "internal", http.StatusInternalServerError)
	}
}
