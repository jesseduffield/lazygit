package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jesseduffield/lazygit/pkg/integration/clients"
)

var usage = `
Usage:
	See https://github.com/jesseduffield/lazygit/tree/master/pkg/integration/README.md

	CLI mode:
		> go run cmd/integration_test/main.go cli [--slow] [--sandbox] <test1> <test2> ...
	If you pass no test names, it runs all tests
	Accepted environment variables:
	KEY_PRESS_DELAY (e.g. 200): the number of milliseconds to wait between keypresses
	MODE:
		* ask (default): if a snapshot test fails, asks if you want to update the snapshot
		* check: if a snapshot test fails, exits with an error
		* update: if a snapshot test fails, updates the snapshot
		* sandbox: uses the test's setup step to run the test in a sandbox where you can do whatever you want

	TUI mode:
		> go run cmd/integration_test/main.go tui
	This will open up a terminal UI where you can run tests

	Help:
		> go run cmd/integration_test/main.go help
`

func main() {
	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	switch os.Args[1] {
	case "help":
		fmt.Println(usage)
	case "cli":
		testNames := os.Args[2:]
		slow := false
		sandbox := false
		// get the next arg if it's --slow
		if len(os.Args) > 2 {
			if os.Args[2] == "--slow" || os.Args[2] == "-slow" {
				testNames = os.Args[3:]
				slow = true
			} else if os.Args[2] == "--sandbox" || os.Args[2] == "-sandbox" {
				testNames = os.Args[3:]
				sandbox = true
			}
		}

		clients.RunCLI(testNames, slow, sandbox)
	case "tui":
		clients.RunTUI()
	default:
		log.Fatal(usage)
	}
}
