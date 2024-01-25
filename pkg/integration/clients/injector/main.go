package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/mitchellh/go-ps"
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

	if os.Getenv(components.WAIT_FOR_DEBUGGER_ENV_VAR) != "" {
		println("Waiting for debugger to attach...")
		for !isDebuggerAttached() {
			time.Sleep(time.Millisecond * 100)
		}

		println("Debugger attached, continuing")
	}

	app.Start(dummyBuildInfo, integrationTest)
}

func getIntegrationTest() integrationTypes.IntegrationTest {
	if daemon.InDaemonMode() {
		// if we've invoked lazygit as a daemon from within lazygit,
		// we don't want to pass a test to the rest of the code.
		return nil
	}

	integrationTestName := os.Getenv(components.TEST_NAME_ENV_VAR)
	if integrationTestName == "" {
		panic(fmt.Sprintf(
			"expected %s environment variable to be set, given that we're running an integration test",
			components.TEST_NAME_ENV_VAR,
		))
	}

	lazygitRootDir := os.Getenv(components.LAZYGIT_ROOT_DIR)
	allTests := tests.GetTests(lazygitRootDir)
	for _, candidateTest := range allTests {
		if candidateTest.Name() == integrationTestName {
			return candidateTest
		}
	}

	panic("Could not find integration test with name: " + integrationTestName)
}

// Returns whether we are running under a debugger. It uses a heuristic to find
// out: when using dlv, it starts a debugserver executable (which is part of
// lldb), and the debuggee becomes a child process of that. So if the name of
// our parent process is "debugserver", we run under a debugger. This works even
// if the parent process used to be the shell and you then attach to the running
// executable.
//
// On Mac this works with VS Code, with the Jetbrains Goland IDE, and when using
// dlv attach in a terminal. I have not been able to verify that it works on
// other platforms, it may have to be adapted there.
func isDebuggerAttached() bool {
	process, err := ps.FindProcess(os.Getppid())
	if err != nil {
		return false
	}
	return process.Executable() == "debugserver"
}
