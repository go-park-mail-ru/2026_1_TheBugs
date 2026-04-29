package main

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/app"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/logger"
)

func main() {
	log := logger.New(string(entity.AuthService))
	err := config.Read(log)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.RunAuthService(&config.Config, log)

}
