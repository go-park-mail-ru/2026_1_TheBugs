package main

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/app"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	launchLogger := logrus.New()
	err := config.Read(launchLogger)
	if err != nil {
		launchLogger.Fatalf("Config error: %s", err)
	}
	appLogger := logger.New(&config.Config)

	app.Run(&config.Config, appLogger)
}
