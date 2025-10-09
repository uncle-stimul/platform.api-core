package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(level string, format string) *logrus.Logger {
	logger := logrus.New()

	logger.Out = os.Stdout

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	if format == "json" {
		logger.SetFormatter(&CustomJSONFormat{})
	} else {
		logger.SetFormatter(&CustomStringFormat{})
	}

	return logger
}
