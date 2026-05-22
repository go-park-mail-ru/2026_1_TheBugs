package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/smtp"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5"
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
	senderRepo := smtp.NewSMTPSender(config.Config.SMTP.Host, config.Config.SMTP.Port, config.Config.SMTP.Email, config.Config.SMTP.Pwd)
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sql := `UPDATE posters_promotions SET status = 'expiered' WHERE ends_at < NOW()`

		ct, err := pool.Exec(ctx, sql)
		if err != nil {
			logger.WithError(err).Error("Failed check promotions")
			return
		}

		if ct.RowsAffected() > 0 {
			logger.Infof("Expiered: %d promotions", ct.RowsAffected())
		}
	})
	c.AddFunc("0 10 * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sql := `SELECT u.email 
        FROM posters_promotions p 
        JOIN users u ON p.user_id = u.id
        WHERE p.status = 'active' 
		  AND p.is_notification_sent = false
          AND p.ends_at > NOW() 
          AND p.ends_at < NOW() + INTERVAL '16 hours'`

		rows, err := pool.Query(ctx, sql)
		if err != nil {
			logger.WithError(err).Error("Failed check expire promotions")
			return
		}
		defer rows.Close()
		emails, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			logger.WithError(err).Error("Failed to collect emails")
			return
		}
		for _, email := range emails {
			err := senderRepo.SendPromotionExpier(ctx, email)
			if err == nil {
				_, err := pool.Exec(ctx, "UPDATE posters_promotions p SET is_notification_sent = true FROM users u WHERE p.user_id = u.id AND u.email = $1 AND p.status='active' ", email)
				if err != nil {
					logger.WithError(err).Error("Failed to update")
				}
			}
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
