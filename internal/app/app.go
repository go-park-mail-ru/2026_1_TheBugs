package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi"
	authHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	complexHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/complex"
	posterHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	userHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/elasticsearch"
	minioRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/minio"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/openrouter"
	tokensRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/redis/tokens"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/smtp"
	uowSql "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/uow"
	authUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	complexUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
	posterUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	userUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/user"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

func Run(cfg *config.ProjectConfig, logger *logrus.Logger) {
	ai := openrouter.New(config.Config.OpenRouter.APIKey, config.Config.OpenRouter.Model)
	dsn := dsn.BuildDSN(cfg.Postgres)
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port), Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	escfg := es.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%d", cfg.ES.Host, cfg.ES.Port),
		},
		Username: "foo",
		Password: "bar",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 5,
			DialContext:           (&net.Dialer{Timeout: time.Second * 5}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	esClient, err := es.NewClient(escfg)
	if err != nil {
		log.Fatalf("es.NewClient: %v", err)
	}
	esRepo := elasticsearch.NewESRepo(esClient)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}
	senderRepo := smtp.NewSMTPSender(config.Config.SMTP.Host, config.Config.SMTP.Port, config.Config.SMTP.Email, config.Config.SMTP.Pwd)

	uow := uowSql.NewSQLStorage(pool)
	tokenRepo := tokensRepo.NewTokenRepo(rdb)
	fileRepo, err := minioRepo.NewFileRepo(minioClient, cfg.Bucket)
	if err != nil {
		log.Fatalf("cannot create file repo: %v", err)
	}

	posterUC := posterUC.NewPosterUseCase(uow, fileRepo, esRepo, ai)
	posterHandler := posterHandler.NewPosterHandler(posterUC)

	authUC := authUC.NewAuthUseCase(uow, tokenRepo, senderRepo)
	authHandler := authHandler.NewAuthHandler(authUC)

	UtilityCompanyUC := complexUC.NewUtilityCompanyUseCase(uow.UtilityCompany())
	utilityCompanyHandler := complexHandler.NewUtilityCompanyHandler(UtilityCompanyUC)

	userUC := userUC.NewUserUseCase(uow, fileRepo)
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
