package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/domains"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			seed := rand.New(rand.NewSource(time.Now().UnixNano()))
			reqId := fmt.Sprintf("%016x", seed.Int())[:10]
			ctx := context.WithValue(r.Context(), domains.ReqID{}, reqId)

			midLogger := logger.WithFields(logrus.Fields{
				"request_id":  reqId,
				"method":      r.Method,
				"remote_addr": r.RemoteAddr,
				"path":        r.URL.Path,
			})

			ctxLog := logger.WithField("request_id", reqId)
			ctx = ctxLogger.SetLogger(ctx, ctxLog)

			midLogger.Info("request started")
			startTime := time.Now()
			defer func() {
				duration := time.Since(startTime)
				midLogger.WithField("duration", duration).Info("request end")
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
