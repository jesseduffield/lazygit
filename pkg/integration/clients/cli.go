package clients

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/generics/slices"
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

func RunCLI(testNames []string, slow bool, sandbox bool) {
	keyPressDelay := tryConvert(os.Getenv("KEY_PRESS_DELAY"), 0)
	if slow {
		keyPressDelay = SLOW_KEY_PRESS_DELAY
	}

	var mode components.Mode
	if sandbox {
		mode = components.SANDBOX
	} else {
		mode = getModeFromEnv()
	}

	err := components.RunTests(
		getTestsToRun(testNames),
		log.Printf,
		runCmdInTerminal,
		runAndPrintFatalError,
		mode,
		keyPressDelay,
		1,
	)
	if err != nil {
		log.Print(err.Error())
	}
}

func runAndPrintFatalError(test *components.IntegrationTest, f func() error) {
	if err := f(); err != nil {
		log.Fatalf(err.Error())
	}
}

func getTestsToRun(testNames []string) []*components.IntegrationTest {
	allIntegrationTests := tests.GetTests()
	var testsToRun []*components.IntegrationTest

	if len(testNames) == 0 {
		return allIntegrationTests
	}

	testNames = slices.Map(testNames, func(name string) string {
		// allowing full test paths to be passed for convenience
		return strings.TrimSuffix(
			regexp.MustCompile(`.*pkg/integration/tests/`).ReplaceAllString(name, ""),
			".go",
		)
	})

outer:
	for _, testName := range testNames {
		// check if our given test name actually exists
		for _, test := range allIntegrationTests {
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
	case "update":
		return components.UPDATE_SNAPSHOT
	case "sandbox":
		return components.SANDBOX
	default:
		log.Fatalf("unknown test mode: %s, must be one of [ask, check, update, sandbox]", os.Getenv("MODE"))
		panic("unreachable")
	}
}

func tryConvert(numStr string, defaultVal int) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return defaultVal
	}

	return num
}
