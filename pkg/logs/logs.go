package logs

import (
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// It's important that this package does not depend on any other package because we
// may want to import it from anywhere, and we don't want to create a circular dependency
// (because Go refuses to compile circular dependencies).

// Global is a global logger that can be used anywhere in the app, for
// _development purposes only_. I want to avoid global variables when possible,
// so if you want to log something that's printed when the -debug flag is set,
// you'll need to ensure the struct you're working with has a logger field (
// and most of them do).
// Global is only available if the LAZYGIT_LOG_PATH environment variable is set.
var Global *logrus.Entry

func init() {
	logPath := os.Getenv("LAZYGIT_LOG_PATH")
	if logPath != "" {
		Global = NewDevelopmentLogger(logPath)
	}
}

func NewProductionLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = io.Discard
	logger.SetLevel(logrus.ErrorLevel)
	return formatted(logger)
}

// Separates one run's log entries from the previous run's. Only the first
// logger of a run writes it: with LAZYGIT_LOG_PATH set there are two of them
// for the same file, the global one and the app's.
var runSeparator sync.Once

func NewDevelopmentLogger(logPath string) *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(getLogLevel())

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Unable to log to log file: %v", err)
	}
	runSeparator.Do(func() {
		_, _ = file.WriteString("\n")
	})
	logger.SetOutput(file)
	return formatted(logger)
}

func formatted(log *logrus.Logger) *logrus.Entry {
	// highly recommended: tail -f development.log | humanlog
	// https://github.com/aybabtme/humanlog
	log.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}

	return log.WithFields(logrus.Fields{})
}

func getLogLevel() logrus.Level {
	strLevel := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(strLevel)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}
