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
	posterrepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/poster"
	userrepo "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/user"
	authuc "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	posteruc "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"

	"github.com/gorilla/mux"
)

//go:generate swag init -g ./internal/app/app.go --parseInternal -o ./docs

// @title Cian
// @version 1.0
// @description Backend для сервиса Циан, команды The Bugs
// @host localhost:8000
// @BasePath /api

func Run(cfg *config.ProjectConfig) {
	posterRepo := posterrepo.NewPosterRepo()
	posterUC := posteruc.NewPosterUseCase(posterRepo)
	posterHandler := posterHandler.NewPosterHandler(posterUC)

	userRepo := userrepo.NewUserRepo()
	userUC := authuc.NewAuthUseCase(userRepo)
	userHandler := authHandler.NewAuthHandler(userUC)

	r := mux.NewRouter()

	restapi.RegisterHandlers(r, userHandler, posterHandler)
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
