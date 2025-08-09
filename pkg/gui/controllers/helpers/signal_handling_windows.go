package helpers

import (
	"github.com/sirupsen/logrus"
)

func canSuspendApp() bool {
	return false
}

func sendStopSignal() error {
	return nil
}

func installResumeSignalHandler(log *logrus.Entry, onResume func() error) {
}
