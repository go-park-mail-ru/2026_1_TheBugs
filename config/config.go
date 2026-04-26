package config

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Config ProjectConfig
var JWTKeys RSAKeys

var DevCors = CORS{
	AllowedHosts: []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:80", "http://localhost", "http://localhost:8000", "http://127.0.0.1:5500"},
	CookieHost:   "localhost",
	URL:          "http://localhost:8000",
}
var ProdCors = CORS{
	AllowedHosts: []string{"http://dom-deli.ru:80", "http://dom-deli.ru", "https://dom-deli.ru"},
	CookieHost:   "dom-deli.ru",
	URL:          "https://dom-deli.ru",
}

type (
	ProjectConfig struct {
		AppEnv string `yaml:"app-env" env:"APP_ENV" env-default:"dev"`
		CORS
		Server     `yaml:"server"`
		Postgres   `yaml:"postgres"`
		Redis      `yaml:"redis"`
		JWT        `yaml:"jwt"`
		OAuth      `yaml:"oauth"`
		SMTP       `yaml:"smtp"`
		Minio      `yaml:"minio"`
		ES         `yaml:"es"`
		OpenRouter `yaml:"openrouter"`
	}
	Redis struct {
		Host     string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
		Port     int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
		DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
	}
	CORS struct {
		AllowedHosts []string
		CookieHost   string
		URL          string
	}
	OAuth struct {
		VKClientID         string `yaml:"vk-client-id"    env:"VKClientID" env-default:"client_id"`
		YandexClientID     string `yaml:"yandex-client-id"    env:"YandexClientID" env-default:"client_id"`
		VKRedirectURI      string `yaml:"vk-redirect-uri"    env:"VKRedirectURI" env-default:"https://dom-deli.ru/oauth/vk"`
		YandexRedirectURI  string `yaml:"yandex-redirect-uri"    env:"YandexRedirectURI" env-default:"https://dom-deli.ru/oauth/yandex"`
		YandexClientSecret string `yaml:"yandex-client-secret"    env:"YandexClientSecret" env-default:"client_secret"`
	}
	SMTP struct {
		Host  string `yaml:"host" env:"SMTP_HOST" env-default:"localhost"`
		Port  int    `yaml:"port" env:"SMTP_PORT" env-default:"1025"`
		Email string `yaml:"email" env:"SMTP_EMAIL" env-default:"admin@dom-deli.ru"`
		Pwd   string `yaml:"pwd" env:"SMTP_PWD" env-default:"1025"`
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
		RecoverExp       time.Duration `yaml:"recover-exp"    env:"JWT_RECOVER_EXP" env-default:"5m"`
	}
	RSAKeys struct {
		PublicKey  *rsa.PublicKey
		PrivateKey *rsa.PrivateKey
	}

	Minio struct {
		Endpoint   string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
		AccessKey  string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-default:"admin123"`
		SecretKey  string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-default:"admin123"`
		Bucket     string `yaml:"bucket" env:"MINIO_BUCKET" env-default:"media"`
		PublicHost string `yaml:"public_host" env:"MINIO_PUBLIC_HOST" env-default:"http://localhost:9000"`
	}

	ES struct {
		Host string `yaml:"host" env:"ES_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"ES_PORT" env-default:"9200"`
	}
	OpenRouter struct {
		APIKey string `yaml:"api_key" env:"OPENROUTER_API_KEY" env-default:""`
		Model  string `yaml:"model" env:"OPENROUTER_MODEL" env-default:"nvidia/nemotron-3-nano-30b-a3b:free"`
	}
)

func Read(log *logrus.Logger) error {
	var err error
	// c:/Users/Артемий/OneDrive/Desktop/code/2026_1_TheBugs/ этот оставил для дебага у вас будет свой
	if err := godotenv.Load(".env"); err != nil {
		log.Warnf("No .env file found: %v", err)
	}
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
