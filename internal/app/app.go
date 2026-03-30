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
	minioRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/minio"
	tokensRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/redis/tokens"
	uowSql "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/uow"
	authUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	complexUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
	posterUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

// переименовать endpoint
func Run(cfg *config.ProjectConfig) {
	dsn := dsn.BuildDSN(cfg.Postgres)
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port), Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	logger := logrus.New()

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}

	uow := uowSql.NewSQLStorage(pool)
	tokenRepo := tokensRepo.NewTokenRepo(rdb)
	fileRepo, err := minioRepo.NewFileRepo(minioClient, cfg.Bucket)
	if err != nil {
		log.Fatalf("cannot create file repo: %v", err)
	}

	posterUC := posterUC.NewPosterUseCase(uow, fileRepo)
	posterHandler := posterHandler.NewPosterHandler(posterUC)

	authUC := authUC.NewAuthUseCase(uow, tokenRepo)
	authHandler := authHandler.NewAuthHandler(authUC)

	UtilityCompanyUC := complexUC.NewUtilityCompanyUseCase(uow.UtilityCompany())
	UtilityCompanyHandler := complexHandler.NewUtilityCompanyHandler(UtilityCompanyUC)

	r := mux.NewRouter()
	restapi.RegisterHandlers(r, logger, authHandler, posterHandler, UtilityCompanyHandler)
	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         serverAddress,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Printf("start listen: %s \n", serverAddress)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)

}
