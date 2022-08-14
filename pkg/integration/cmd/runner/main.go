package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
)

// see pkg/integration/README.md

// The purpose of this program is to run integration tests. It does this by
// building our injector program (in the sibling injector directory) and then for
// each test we're running, invoke the injector program with the test's name as
// an environment variable. Then the injector finds the test and passes it to
// the lazygit startup code.

// If invoked directly, you can specify tests to run by passing their names as positional arguments

func main() {
	err := components.RunTests(
		getTestsToRun(),
		log.Printf,
		runCmdInTerminal,
		func(test *components.IntegrationTest, f func() error) {
			if err := f(); err != nil {
				log.Print(err.Error())
			}
		},
		getModeFromEnv(),
	)
	if err != nil {
		log.Print(err.Error())
	}
}

func getTestsToRun() []*components.IntegrationTest {
	var testsToRun []*components.IntegrationTest

	if len(os.Args) < 2 {
		return tests.Tests
	}

outer:
	for _, testName := range os.Args[1:] {
		// check if our given test name actually exists
		for _, test := range tests.Tests {
			if test.Name() == testName {
				testsToRun = append(testsToRun, test)
				continue outer
			}
		}
		log.Fatalf("test %s not found. Perhaps you forgot to add it to `pkg/integration/integration_tests/tests.go`?", testName)
	}

	return testsToRun
}

func runCmdInTerminal(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getModeFromEnv() components.Mode {
	switch os.Getenv("MODE") {
	case "", "ask":
		return components.ASK_TO_UPDATE_SNAPSHOT
	case "check":
		return components.CHECK_SNAPSHOT
	case "updateSnapshot":
		return components.UPDATE_SNAPSHOT
	case "sandbox":
		return components.SANDBOX
	default:
		log.Fatalf("unknown test mode: %s, must be one of [test, record, updateSnapshot, sandbox]", os.Getenv("MODE"))
		panic("unreachable")
	}
}
