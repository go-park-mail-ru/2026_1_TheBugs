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
	orderHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/order"
	posterHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	promHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/promotion"
	userHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/user"
	minioRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/minio"
	uowSql "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/uow"
	orderUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/order"
	promUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/promotion"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/mux"
)

func Run(cfg *config.ProjectConfig, logger *logrus.Logger) {
	dsn := dsn.BuildDSN(cfg.Postgres)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}
	uow := uowSql.NewSQLStorage(pool)
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("cannot create minio client: %v", err)
	}

	fileRepo, err := minioRepo.NewFileRepo(minioClient, cfg.Bucket)
	if err != nil {
		log.Fatalf("cannot create file repo: %v", err)
	}

	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.AuthService.Host, cfg.AuthService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("cannot dial grpc server: %v", err)
	}
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.UserService.Host, cfg.UserService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("cannot dial grpc server: %v", err)
	}
	posterConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.PosterService.Host, cfg.PosterService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("cannot dial poster grpc server: %v", err)
	}
	complexConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.ComplexService.Host, cfg.ComplexService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("cannot dial grpc server: %v", err)
	}
	authHandler := authHandler.NewAuthHandler(authConn)
	userHandler := userHandler.NewUserHandler(userConn)
	posterHandler := posterHandler.NewPosterHandler(posterConn)
	utilityCompanyHandler := complexHandler.NewUtilityCompanyHandler(complexConn)

	orderUC := orderUC.NewOrderUseCase(uow, fileRepo)
	orderHandler := orderHandler.NewOrderHandler(orderUC)

	promotionUC := promUC.NewPromotionUseCase(uow)
	promotionHandler := promHandler.NewPromotionHandler(promotionUC)

	r := mux.NewRouter()
	restapi.RegisterHandlers(r, logger, authHandler, posterHandler, utilityCompanyHandler, userHandler, orderHandler, promotionHandler)
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
