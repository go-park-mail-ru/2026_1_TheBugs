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
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/mux"
)

func Run(cfg *config.ProjectConfig, logger *logrus.Logger) {

	// rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port), Password: cfg.Redis.Password, DB: cfg.Redis.DB})

	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.AuthService.Host, cfg.AuthService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot dial grpc server: %v", err)
	}
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.UserService.Host, cfg.UserService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot dial grpc server: %v", err)
	}
	posterConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.PosterService.Host, cfg.PosterService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot dial poster grpc server: %v", err)
	}
	complexConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.ComplexService.Host, cfg.ComplexService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot dial grpc server: %v", err)
	}
	authHandler := authHandler.NewAuthHandler(authConn)
	userHandler := userHandler.NewUserHandler(userConn)
	utilityCompanyHandler := complexHandler.NewUtilityCompanyHandler(complexConn)
	posterHandler := posterHandler.NewPosterHandler(posterConn)

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
