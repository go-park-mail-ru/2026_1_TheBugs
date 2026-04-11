package middleware

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "OPTIONS" || r.Method == "HEADERS" {
			next.ServeHTTP(w, r)
			return
		}
		cookieToken, err := r.Cookie("csrf_token")
		if err != nil {
			utils.WriteError(w, "Not Found CSRF", 403)
			return
		}
		headerToken := r.Header.Get("X-CSRF-TOKEN")
		if cookieToken == nil || headerToken == "" || cookieToken.Value != headerToken {
			utils.WriteError(w, "Not Found CSRF", 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := generateCSRFToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})
	log.Println(token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"csrf_token": token})
}
