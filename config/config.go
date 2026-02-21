package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var Config ProjectConfig

type (
	ProjectConfig struct {
		Server   `yaml:"server"`
		Postgres `yaml:"postgres"`
	}

	Server struct {
		Host         string        `yaml:"host"    env:"SRV_HOST" env-default:"localhost"`
		Port         int           `yaml:"port"    env:"SRV_PORT" env-default:"8000"`
		WriteTimeout time.Duration `yaml:"write-timeout"    env:"SRV_WRITE_TM" env-default:"5s"`
		ReadTimeout  time.Duration `yaml:"read-timeout"    env:"SRV_READ_TM" env-default:"5s"`
		IdleTimeout  time.Duration `yaml:"idle-timeout"    env:"SRV_IDLE_TM" env-default:"20s"`
	}

	Postgres struct {
		Host            string        `yaml:"host"    env:"PG_HOST" env-default:"localhost"`
		Port            int           `yaml:"port" env:"PG_PORT" env-default:"5432"`
		User            string        `yaml:"user" env:"PG_USER" env-default:"thebugs"`
		Password        string        `yaml:"password" env:"PG_PASSWORD" env-default:"thebugs"`
		Database        string        `yaml:"database" env:"PG_DB" env-default:"main"`
		SslMode         string        `yaml:"sslmode" env:"PG_SSL_MODE" env-default:"disable"`
		MaxOpenConns    int           `yaml:"max-open-connections" env:"PG_MAX_OPEN_CONN" env-default:"10"`
		ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime" env:"PG_CONN_MAX_LIFETIME" env-default:"30s"`
	}
)

func Read() error {
	if err := cleanenv.ReadConfig("config/config.yaml", &Config); err != nil {
		return fmt.Errorf("error while reading application configuration: %w", err)
	}

	if err := cleanenv.ReadEnv(&Config); err != nil {
		return fmt.Errorf("error creating configuration object: %w", err)
	}
	log.Println("reading configuration is successful")
	return nil
}
