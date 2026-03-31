package main

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	err := config.Read(logger)
	if err != nil {
		logger.Fatalf("Config error: %s", err)
	}

	app.Run(&config.Config, logger)
}
