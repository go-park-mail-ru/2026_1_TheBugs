package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	complexHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/complex"
	complexGRPC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/interceptor"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	prom "github.com/go-park-mail-ru/2026_1_TheBugs/internal/metrics"
	complexRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/complex"
	complexUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gorilla/mux"
)

func RunComplexService(cfg *config.ProjectConfig, logger *logrus.Logger) {
	dsn := dsn.BuildDSN(cfg.Postgres)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}
	repo := complexRepo.NewUtilityCompanyRepo(pool)

	app := mux.NewRouter()
	uc := complexUC.NewUtilityCompanyUseCase(repo)

	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcAsddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.UserService.Port)

	metrics := prom.NewMetricsMiddleware()
	metrics.Register(entity.ServiceType(cfg.ServiceName))

	app.Handle("/metrics", promhttp.Handler())
	app.Use(middleware.LoggingMiddleware(logger))

	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      app,
		Addr:         serverAddress,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	logger.Infof("start http listen: %s", serverAddress)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingServerInterceptor(logger),
			metrics.MetricsGRPCInterceptor,
		),
	)
	complexGRPC.RegisterComplexServiceServer(
		grpcServer,
		complexHandler.NewComplexServiceServer(uc),
	)
	logger.Info("start grpc listen: %s", grpcAsddr)
	go func() {
		lis, err := net.Listen("tcp", grpcAsddr)
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
	logger.Warn("shutting http down")
	grpcServer.GracefulStop()
	logger.Warn("shutting grpc down")
	os.Exit(0)
}
