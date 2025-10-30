package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetLevel(logrus.DebugLevel)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Info(msg string) {
	logger.Info(msg)
}

func Warningf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Warning(msg string) {
	logger.Warn(msg)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Writer() *io.PipeWriter {
	return logger.Writer()
}
