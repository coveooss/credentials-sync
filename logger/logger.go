package logger

import (
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func initSentry() *logrus.Logger {
	newLogger := logrus.New()

	sentryDsn, ok := os.LookupEnv("SENTRY_DSN")
	if !ok {
		newLogger.Info("Not using sentry, SENTRY_DSN not set")
		return newLogger
	}
	for _, sentryVariable := range []string{"SENTRY_ENVIRONMENT", "SENTRY_RELEASE"} {
		if sentryVariableValue := os.Getenv(sentryVariable); sentryVariableValue == "" {
			newLogger.Infof("Not using sentry, %s not set", sentryVariable)
			return newLogger
		}
	}

	err := sentry.Init(sentry.ClientOptions{
		Debug: false,
	})
	if err != nil {
		logrus.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	hook, err := logrus_sentry.NewSentryHook(sentryDsn, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	if err == nil {
		hook.Timeout = 2 * time.Second
		newLogger.AddHook(hook)
	}

	newLogger.Info("Sentry initialized")
	return newLogger
}

var Log *logrus.Logger = initSentry()