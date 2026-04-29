package auth

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/app"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/logger"
)

func main() {
	logger := logger.New()
	err := config.Read(logger)
	if err != nil {
		logger.Fatalf("Config error: %s", err)
	}

	app.RunAuthService(&config.Config, logger)

}
