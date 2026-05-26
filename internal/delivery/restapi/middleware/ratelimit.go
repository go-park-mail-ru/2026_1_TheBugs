package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
)

func RateLimitMiddleware(rateLimitUC delivery.RateLimitUseCase) func(limit int, interval time.Duration) func(next http.Handler) http.Handler {
	return func(limit int, interval time.Duration) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log := ctxLogger.GetLogger(r.Context()).WithField("op", "RateLimitMiddleware")

				allowed, err := rateLimitUC.CheckIPLimit(r.Context(), r.RemoteAddr, limit, interval)
				if err != nil {
					log.Errorf("rateLimitUC.CheckIPLimit: %s", err)
					utils.HandelError(w, err)
					return
				}
				if !allowed {
					log.Warnf("rate limit exceeded for IP %s (limit=%d)", r.RemoteAddr, limit)
					utils.WriteError(w, "too many requests", http.StatusTooManyRequests)
					return
				}

				next.ServeHTTP(w, r)
			})
		}
	}
}
