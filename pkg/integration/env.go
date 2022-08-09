package integration

import (
	"os"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/integration/integration_tests"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

// NEW integration test format stuff

func IntegrationTestName() string {
	return os.Getenv("LAZYGIT_TEST_NAME")
}

func PlayingIntegrationTest() bool {
	return IntegrationTestName() != ""
}

func CurrentIntegrationTest() (types.Test, bool) {
	if !PlayingIntegrationTest() {
		return nil, false
	}

	return slices.Find(integration_tests.Tests, func(test types.Test) bool {
		return test.Name() == IntegrationTestName()
	})
}
