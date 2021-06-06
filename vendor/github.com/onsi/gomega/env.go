package gomega

import (
	"os"

	"github.com/onsi/gomega/internal/defaults"
)

const (
	ConsistentlyDurationEnvVarName        = "GOMEGA_DEFAULT_CONSISTENTLY_DURATION"
	ConsistentlyPollingIntervalEnvVarName = "GOMEGA_DEFAULT_CONSISTENTLY_POLLING_INTERVAL"
	EventuallyTimeoutEnvVarName           = "GOMEGA_DEFAULT_EVENTUALLY_TIMEOUT"
	EventuallyPollingIntervalEnvVarName   = "GOMEGA_DEFAULT_EVENTUALLY_POLLING_INTERVAL"
)

func init() {
	defaults.SetDurationFromEnv(
		os.Getenv,
		SetDefaultConsistentlyDuration,
		ConsistentlyDurationEnvVarName,
	)

	defaults.SetDurationFromEnv(
		os.Getenv,
		SetDefaultConsistentlyPollingInterval,
		ConsistentlyPollingIntervalEnvVarName,
	)

	defaults.SetDurationFromEnv(
		os.Getenv,
		SetDefaultEventuallyTimeout,
		EventuallyTimeoutEnvVarName,
	)

	defaults.SetDurationFromEnv(
		os.Getenv,
		SetDefaultEventuallyPollingInterval,
		EventuallyPollingIntervalEnvVarName,
	)
}
