package main

import (
	"os"

	"github.com/jesseduffield/lazygit/pkg/app"
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
	integrationTestName := os.Getenv("LAZYGIT_TEST_NAME")
	if integrationTestName == "" {
		panic("expected LAZYGIT_TEST_NAME environment variable to be set, given that we're running an integration test")
	}

	// unsetting so that if we run lazygit in as a 'daemon' we don't think we're trying to run a test again
	os.Unsetenv("LAZYGIT_TEST_NAME")
	for _, candidateTest := range integration.Tests {
		if candidateTest.Name() == integrationTestName {
			return candidateTest
		}
	}

	panic("Could not find integration test with name: " + integrationTestName)
}
