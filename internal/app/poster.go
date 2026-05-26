package app

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	posterGRPC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/interceptor"
	posterHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	prom "github.com/go-park-mail-ru/2026_1_TheBugs/internal/metrics"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/elasticsearch"
	minioRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/minio"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/openrouter"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/redis/cache"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/osm"
	uowSql "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/uow"
	posterUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RunPosterService(cfg *config.ProjectConfig, logger *logrus.Logger) {
	pool, err := dsn.OpenDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port), Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	cacheRepo := cache.NewRedisCache(rdb)

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

	esCfg := es.Config{
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

	esClient, err := es.NewClient(esCfg)
	if err != nil {
		log.Fatalf("es.NewClient: %v", err)
	}

	esRepo := elasticsearch.NewESRepo(esClient)
	ai := openrouter.New(config.Config.OpenRouter.APIKey, config.Config.OpenRouter.Model)
	maps := osm.NewOSMRepo()

	posterUC := posterUC.NewPosterUseCase(uow, fileRepo, esRepo, ai, cacheRepo, maps)

	app := mux.NewRouter()

	metrics := prom.NewMetricsMiddleware()
	metrics.Register(entity.ServiceType(cfg.ServiceName))

	app.Handle("/metrics", promhttp.Handler())
	app.Use(middleware.LoggingMiddleware(logger))

	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.PosterService.Port)

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

	posterGRPC.RegisterPosterServiceServer(
		grpcServer,
		posterHandler.NewPosterServiceServer(posterUC),
	)

	logger.Infof("start grpc listen: %s", grpcAddr)
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
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
