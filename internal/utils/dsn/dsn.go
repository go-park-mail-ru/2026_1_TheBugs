package dsn

import (
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
)

func BuildDSN(cfg config.Postgres) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SslMode,
	)
}
