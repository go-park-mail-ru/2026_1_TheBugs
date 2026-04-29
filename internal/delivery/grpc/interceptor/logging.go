package interceptor

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func LoggingServerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		seed := rand.New(rand.NewSource(time.Now().UnixNano()))
		reqId := fmt.Sprintf("%016x", seed.Int())[:10]
		ctx = context.WithValue(ctx, entity.ReqID{}, reqId)

		midLogger := logger.WithFields(logrus.Fields{
			"request_id":  reqId,
			"method":      info.FullMethod,
			"remote_addr": info.Server,
		})

		log := logger.WithField("request_id", reqId)
		ctx = ctxLogger.SetLogger(ctx, log)

		midLogger.Info("request started")
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			midLogger.WithField("duration", duration).Info("request end")
		}()

		resp, err := handler(ctx, req)

		return resp, err
	}
}
