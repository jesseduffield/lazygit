package app

import (
	"io"
	"log"
	"os"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/sirupsen/logrus"
)

func newLogger(config config.AppConfigurer) *logrus.Entry {
	var log *logrus.Logger
	if config.GetDebug() {
		log = newDevelopmentLogger()
	} else {
		log = newProductionLogger()
	}

	// highly recommended: tail -f development.log | humanlog
	// https://github.com/aybabtme/humanlog
	log.Formatter = &logrus.JSONFormatter{}

	return log.WithFields(logrus.Fields{})
}

func newProductionLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = io.Discard
	log.SetLevel(logrus.ErrorLevel)
	return log
}

func newDevelopmentLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(getLogLevel())
	logPath, err := config.LogPath()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Unable to log to log file: %v", err)
	}
	logger.SetOutput(file)
	return logger
}

func getLogLevel() logrus.Level {
	strLevel := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(strLevel)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}
