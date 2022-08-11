package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/integration"
	"github.com/jesseduffield/lazygit/pkg/integration/helpers"
)

// see pkg/integration/README.md

// If invoked directly, you can specify tests to run by passing them as positional arguments.

func main() {
	mode := integration.GetModeFromEnv()
	includeSkipped := os.Getenv("INCLUDE_SKIPPED") == "true"
	var testsToRun []*helpers.IntegrationTest

	if len(os.Args) > 1 {
	outer:
		for _, testName := range os.Args[1:] {
			// check if our given test name actually exists
			for _, test := range integration.Tests {
				if test.Name() == testName {
					testsToRun = append(testsToRun, test)
					continue outer
				}
			}
			log.Fatalf("test %s not found. Perhaps you forgot to add it to `pkg/integration/integration_tests/tests.go`?", testName)
		}
	} else {
		testsToRun = integration.Tests
	}

	testNames := slices.Map(testsToRun, func(test *helpers.IntegrationTest) string {
		return test.Name()
	})

	err := integration.RunTests(
		log.Printf,
		runCmdInTerminal,
		func(test *helpers.IntegrationTest, f func() error) {
			if !slices.Contains(testNames, test.Name()) {
				return
			}
			if err := f(); err != nil {
				log.Print(err.Error())
			}
		},
		mode,
		includeSkipped,
	)
	if err != nil {
		log.Print(err.Error())
	}
}

func runCmdInTerminal(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
