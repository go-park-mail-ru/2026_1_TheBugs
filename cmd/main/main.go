package main

import (
	"2026_1_the_bugs/config"
	"2026_1_the_bugs/internal/app"
	"log"
)

func main() {
	err := config.Read()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(&config.Config)
}
