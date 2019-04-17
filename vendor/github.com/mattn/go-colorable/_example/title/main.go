package main

import (
	"fmt"
	"os"
	. "github.com/mattn/go-colorable"
)

func main() {
	out := NewColorableStdout()
	fmt.Fprint(out, "\x1B]0;TITLE Changed\007(See title and hit any key)")
	var c [1]byte
	os.Stdin.Read(c[:])
}
