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
	posterHandler "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	authRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/auth"
	posterrepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/poster"
	userRepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/user"
	authUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	posterUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/dsn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gorilla/mux"
)

func Run(cfg *config.ProjectConfig) {
	dsn := dsn.BuildDSN(cfg.Postgres)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("cannot create pgx pool: %v", err)
	}

	posterRepo := posterrepo.NewPosterRepo(pool)
	posterUC := posterUC.NewPosterUseCase(posterRepo)
	posterHandler := posterHandler.NewPosterHandler(posterUC)

	userRepo := userRepo.NewUserRepo(pool)
	authRepo := authRepo.NewAuthRepo(pool)
	authUC := authUC.NewAuthUseCase(userRepo, authRepo)
	authHandler := authHandler.NewAuthHandler(authUC)

	r := mux.NewRouter()
	restapi.RegisterHandlers(r, authHandler, posterHandler)
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
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)

}
