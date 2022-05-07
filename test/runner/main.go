package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/integration"
	"github.com/stretchr/testify/assert"
)

// see docs/Integration_Tests.md
// This file can be invoked directly, but you might find it easier to go through
// test/lazyintegration/main.go, which provides a convenient gui wrapper to integration tests.
//
// If invoked directly, you can specify a test by passing it as the first argument.
// You can also specify that you want to record a test by passing MODE=record
// as an env var.

func main() {
	mode := integration.GetModeFromEnv()
	speedEnv := os.Getenv("SPEED")
	includeSkipped := os.Getenv("INCLUDE_SKIPPED") == "true"
	selectedTestName := os.Args[1]

	err := integration.RunTests(
		log.Printf,
		runCmdInTerminal,
		func(test *integration.Test, f func(*testing.T) error) {
			if selectedTestName != "" && test.Name != selectedTestName {
				return
			}
			if err := f(nil); err != nil {
				log.Print(err.Error())
			}
		},
		mode,
		speedEnv,
		func(_t *testing.T, expected string, actual string, prefix string) { //nolint:thelper
			assert.Equal(MockTestingT{}, expected, actual, fmt.Sprintf("Unexpected %s. Expected:\n%s\nActual:\n%s\n", prefix, expected, actual))
		},
		includeSkipped,
	)
	if err != nil {
		log.Print(err.Error())
	}
}

type MockTestingT struct{}

func (t MockTestingT) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func runCmdInTerminal(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
