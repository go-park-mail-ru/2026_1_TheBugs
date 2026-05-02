package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	err := config.Read(logger)
	if err != nil {
		logger.Fatalf("Config error: %s", err)
	}
	dsn := dsn.BuildDSN(config.Config.Postgres)
	pool, err := pgxpool.New(context.Background(), dsn)
	c := cron.New()
	c.AddFunc("*/10 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sql := `DELETE FROM users WHERE is_verified = FALSE AND created_at < NOW() - INTERVAL '10 minutes'`

		ct, err := pool.Exec(ctx, sql)
		if err != nil {
			logger.WithError(err).Error("Failed to delete unverified users")
			return
		}

		if ct.RowsAffected() > 0 {
			logger.Infof("Cleanup: deleted %d unverified users", ct.RowsAffected())
		}
	})
	c.Start()
	logger.Info("Cron job started. Press Ctrl+C to stop.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Info("Stopping cron...")
	c.Stop()
	pool.Close()
	logger.Info("Application stopped.")

}
