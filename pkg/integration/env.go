package integration

import (
	"os"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

func Headless() bool {
	return os.Getenv("HEADLESS") != ""
}

// NEW integration test format stuff

func IntegrationTestName() string {
	return os.Getenv("LAZYGIT_TEST_NAME")
}

func CurrentIntegrationTest() (types.Test, bool) {
	if !PlayingIntegrationTest() {
		return nil, false
	}

	return slices.Find(Tests, func(test types.Test) bool {
		return test.Name() == IntegrationTestName()
	})
}

func PlayingIntegrationTest() bool {
	return IntegrationTestName() != ""
}

// OLD integration test format stuff

func Replaying() bool {
	return os.Getenv("REPLAY_EVENTS_FROM") != ""
}

func RecordingEvents() bool {
	return recordEventsTo() != ""
}

func recordEventsTo() string {
	return os.Getenv("RECORD_EVENTS_TO")
}
