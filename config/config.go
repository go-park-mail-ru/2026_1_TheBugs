package config

import (
	"crypto/rsa"
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var Config ProjectConfig
var JWTKeys RSAKeys

var DevCors = CORS{
	AllowedHosts: []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:80", "http://localhost"},
	CookieHost:   "localhost",
}
var ProdCors = CORS{
	AllowedHosts: []string{"http://dom-deli.ru:80", "http://dom-deli.ru", "https://dom-deli.ru"},
	CookieHost:   "dom-deli.ru",
}

type (
	ProjectConfig struct {
		AppEnv string `yaml:"app-env" env:"APP_ENV" env-default:"dev"`
		CORS
		Server   `yaml:"server"`
		Postgres `yaml:"postgres"`
		JWT      `yaml:"jwt"`
	}
	CORS struct {
		AllowedHosts []string
		CookieHost   string
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
	JWT struct {
		PublicKeySource  string        `yaml:"public-key-source"    env:"JWT_PUBLIC_KEY_SOURCE" env-default:"public.pem"`
		PrivateKeySource string        `yaml:"private-key-source"    env:"JWT_PRIVATE_KEY_SOURCE" env-default:"private.pem"`
		AccessExp        time.Duration `yaml:"access-exp"    env:"JWT_ACCESS_EXP" env-default:"15m"`
		RefreshExp       time.Duration `yaml:"refresh-exp"    env:"JWT_REFRESH_EXP" env-default:"24h"`
	}
	RSAKeys struct {
		PublicKey  *rsa.PublicKey
		PrivateKey *rsa.PrivateKey
	}
)

func Read() error {
	var err error
	// c:/Users/Артемий/OneDrive/Desktop/code/2026_1_TheBugs/ этот оставил для дебага у вас будет свой
	if err = cleanenv.ReadConfig("config/config.yaml", &Config); err != nil {
		return fmt.Errorf("error while reading application configuration: %w", err)
	}

	if err = cleanenv.ReadEnv(&Config); err != nil {
		return fmt.Errorf("error creating configuration object: %w", err)
	}
	JWTKeys.PrivateKey, err = LoadPrivateKey(Config.JWT.PrivateKeySource)
	if err != nil {
		return fmt.Errorf("error load private key: %w", err)
	}

	JWTKeys.PublicKey, err = LoadPublicKey(Config.JWT.PublicKeySource)
	if err != nil {
		return fmt.Errorf("error load public key: %w", err)
	}
	if Config.AppEnv == "dev" {
		Config.CORS = DevCors
	} else {
		Config.CORS = ProdCors
	}
	log.Print(Config.CORS.AllowedHosts)
	log.Println("reading configuration is successful")
	return nil
}
