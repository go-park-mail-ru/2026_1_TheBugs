package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}

}
