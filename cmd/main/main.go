package main

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/app"
)

func main() {
	err := config.Read()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(&config.Config)
}
