package logger

import (
	"os"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

func New(cfg *config.ProjectConfig) *logrus.Logger {
	godotenv.Load(".env")
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	lokiURL := cfg.Loki.URL

	if lokiURL != "" {
		opts := lokirus.NewLokiHookOptions().
			WithFormatter(&logrus.JSONFormatter{}).
			WithStaticLabels(lokirus.Labels{
				"service": cfg.ServiceName,
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
		"service": cfg.ServiceName,
		"loki":    lokiURL != "",
	}).Info("logger initialized")

	return log
}
