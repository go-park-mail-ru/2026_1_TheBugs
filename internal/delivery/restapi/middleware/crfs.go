package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
)

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
