package ctxLogger

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/domains"
	"github.com/sirupsen/logrus"
)

func SetLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, domains.LoggerCtx{}, logger)
}

func GetLogger(ctx context.Context) *logrus.Entry {
	log, ok := ctx.Value(domains.LoggerCtx{}).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.New())
	}
	return log
}
