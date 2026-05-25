package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mailru/easyjson"
)

func JSONResponse(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}

}

func EasyJSONResponse(w http.ResponseWriter, status int, obj easyjson.Marshaler) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(obj, w); err != nil {
		log.Printf("easyjson.MarshalToHTTPResponseWriter: %s", err)
	}

}
