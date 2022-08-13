package main

import (
	"fmt"
	"os"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/integration"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

// The purpose of this program is to run lazygit with an integration test passed in.
// We could have done the check on LAZYGIT_TEST_NAME in the root main.go but
// that would mean lazygit would be depending on integration test code which
// would bloat the binary.

// You should not invoke this program directly. Instead you should go through
// pkg/integration/cmd/runner/main.go or pkg/integration/cmd/tui/main.go

func main() {
	dummyBuildInfo := &app.BuildInfo{
		Commit:      "",
		Date:        "",
		Version:     "",
		BuildSource: "integration test",
	}

	integrationTest := getIntegrationTest()

	app.Start(dummyBuildInfo, integrationTest)
}

func getIntegrationTest() integrationTypes.IntegrationTest {
	if daemon.InDaemonMode() {
		// if we've invoked lazygit as a daemon from within lazygit,
		// we don't want to pass a test to the rest of the code.
		return nil
	}

	integrationTestName := os.Getenv(integration.LAZYGIT_TEST_NAME_ENV_VAR)
	if integrationTestName == "" {
		panic(fmt.Sprintf(
			"expected %s environment variable to be set, given that we're running an integration test",
			integration.LAZYGIT_TEST_NAME_ENV_VAR,
		))
	}

	for _, candidateTest := range integration.Tests {
		if candidateTest.Name() == integrationTestName {
			return candidateTest
		}
	}

	panic("Could not find integration test with name: " + integrationTestName)
}
