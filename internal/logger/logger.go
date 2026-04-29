package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

func New(serviceName string) *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(getLogLevel())

	lokiURL := getEnv("LOKI_URL", "http://localhost:3100")

	if lokiURL != "" {
		opts := lokirus.NewLokiHookOptions().
			WithFormatter(&logrus.JSONFormatter{}).
			WithStaticLabels(lokirus.Labels{
				"service": serviceName,
			}).
			WithLevelMap(lokirus.LevelMap{
				logrus.TraceLevel: "trace",
				logrus.DebugLevel: "debug",
				logrus.InfoLevel:  "info",
				logrus.WarnLevel:  "warn",
				logrus.ErrorLevel: "error",
				logrus.FatalLevel: "fatal",
				logrus.PanicLevel: "critical",
			})

		hook := lokirus.NewLokiHookWithOpts(
			lokiURL,
			opts,
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		)

		log.AddHook(hook)
	}

	log.WithFields(logrus.Fields{
		"service": serviceName,
		"loki":    lokiURL != "",
	}).Info("logger initialized")

	return log
}

func getLogLevel() logrus.Level {
	raw := strings.ToLower(getEnv("LOG_LEVEL", "info"))

	level, err := logrus.ParseLevel(raw)
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
