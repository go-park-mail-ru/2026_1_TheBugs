package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi"
	authHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	complexHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/complex"
	posterHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	userHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/user"
	tokensRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/redis/tokens"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/smtp"
	uowSql "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/uow"
	authUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	complexUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
	posterUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	userUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/user"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

func Run(cfg *config.ProjectConfig, logger *logrus.Logger) {
	dsn := dsn.BuildDSN(cfg.Postgres)
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port), Password: cfg.Redis.Password, DB: cfg.Redis.DB})

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}
	senderRepo := smtp.NewSMTPSender(config.Config.SMTP.Host, config.Config.SMTP.Port, config.Config.SMTP.Email, config.Config.SMTP.Pwd)

	uow := uowSql.NewSQLStorage(pool)
	tokenRepo := tokensRepo.NewTokenRepo(rdb)

	posterUC := posterUC.NewPosterUseCase(uow.Posters())
	posterHandler := posterHandler.NewPosterHandler(posterUC)

	authUC := authUC.NewAuthUseCase(uow, tokenRepo, senderRepo)
	authHandler := authHandler.NewAuthHandler(authUC)

	UtilityCompanyUC := complexUC.NewUtilityCompanyUseCase(uow.UtilityCompany())
	utilityCompanyHandler := complexHandler.NewUtilityCompanyHandler(UtilityCompanyUC)

	userUC := userUC.NewUserUseCase(uow)
	userHandler := userHandler.NewUserHandler(userUC)

	r := mux.NewRouter()
	restapi.RegisterHandlers(r, logger, authHandler, posterHandler, utilityCompanyHandler, userHandler)
	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         serverAddress,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	logger.Infof("start listen: %s", serverAddress)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
	logger.Warn("shutting down")
	os.Exit(0)

}
