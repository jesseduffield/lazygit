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
	INPUT_DELAY (e.g. 200): the number of milliseconds to wait between keypresses or mouse clicks

	TUI mode:
		> go run cmd/integration_test/main.go tui
	This will open up a terminal UI where you can run tests

	Help:
		> go run cmd/integration_test/main.go help
`

type flagInfo struct {
	name string // name of the flag; can be used with "-" or "--"
	flag *bool  // a pointer to the variable that should be set to true when this flag is passed
}

// Takes the args that you want to parse (excluding the program name and any
// subcommands), and returns the remaining args with the flags removed
func parseFlags(args []string, flags []flagInfo) []string {
outer:
	for len(args) > 0 {
		for _, f := range flags {
			if args[0] == "-"+f.name || args[0] == "--"+f.name {
				*f.flag = true
				args = args[1:]
				continue outer
			}
		}
		break
	}

	return args
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	switch os.Args[1] {
	case "help":
		fmt.Println(usage)
	case "cli":
		slow := false
		sandbox := false
		waitForDebugger := false
		raceDetector := false
		testNames := parseFlags(os.Args[2:], []flagInfo{
			{"slow", &slow},
			{"sandbox", &sandbox},
			{"debug", &waitForDebugger},
			{"race", &raceDetector},
		})
		clients.RunCLI(testNames, slow, sandbox, waitForDebugger, raceDetector)
	case "tui":
		raceDetector := false
		remainingArgs := parseFlags(os.Args[2:], []flagInfo{
			{"race", &raceDetector},
		})
		if len(remainingArgs) > 0 {
			log.Fatal("tui only supports the -race argument.")
		}
		clients.RunTUI(raceDetector)
	default:
		log.Fatal(usage)
	}
}
