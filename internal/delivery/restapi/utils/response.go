package utils

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}
