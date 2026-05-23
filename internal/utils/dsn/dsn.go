package dsn

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/jackc/pgx/v5/pgxpool"
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

func OpenDB(cfg config.Postgres) (*pgxpool.Pool, error) {
	dsn := BuildDSN(cfg)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConnLifetime = cfg.ConnMaxLifetime
	config.MaxConns = int32(cfg.MaxIdleConns)
	config.MaxConnIdleTime = cfg.ConnMaxIdleTime
	config.ConnConfig.RuntimeParams = map[string]string{
		"statement_timeout": strconv.Itoa(cfg.StatementTimeout),
		"lock_timeout":      strconv.Itoa(cfg.LockTimeout),
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
