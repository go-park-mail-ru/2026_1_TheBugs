package ctxLogger

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/sirupsen/logrus"
)

func SetLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, entity.LoggerCtx{}, logger)
}

func GetLogger(ctx context.Context) *logrus.Entry {
	log, ok := ctx.Value(entity.LoggerCtx{}).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.New())
	}
	return log
}
