package utils

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(obj)
}
