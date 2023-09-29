//go:build ignore

package main

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/cheatsheet"
)

func main() {
	fmt.Printf("Generating cheatsheets in %s...\n", cheatsheet.GetKeybindingsDir())
	cheatsheet.Generate()
}
