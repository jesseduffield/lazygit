package main

import (
	"fmt"
	"os"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

// The purpose of this program is to run lazygit with an integration test passed in.
// We could have done the check on TEST_NAME in the root main.go but
// that would mean lazygit would be depending on integration test code which
// would bloat the binary.

// You should not invoke this program directly. Instead you should go through
// go run cmd/integration_test/main.go

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

	if os.Getenv(components.SANDBOX_ENV_VAR) == "true" {
		// when in sandbox mode we don't want the test controlling the gui
		return nil
	}

	integrationTestName := os.Getenv(components.TEST_NAME_ENV_VAR)
	if integrationTestName == "" {
		panic(fmt.Sprintf(
			"expected %s environment variable to be set, given that we're running an integration test",
			components.TEST_NAME_ENV_VAR,
		))
	}

	allTests := tests.GetTests()
	for _, candidateTest := range allTests {
		if candidateTest.Name() == integrationTestName {
			return candidateTest
		}
	}

	panic("Could not find integration test with name: " + integrationTestName)
}
