package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jesseduffield/lazygit/pkg/cheatsheet"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a command: one of 'generate', 'check'")
	}

	command := os.Args[1]

	switch command {
	case "generate":
		cheatsheet.Generate()
		fmt.Printf("\nGenerated cheatsheets in %s\n", cheatsheet.GetDir())
	case "check":
		cheatsheet.Check()
	default:
		log.Fatal("\nUnknown command. Expected one of 'generate', 'check'")
	}
}
