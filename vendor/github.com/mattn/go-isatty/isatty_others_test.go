// +build !windows

package isatty

import (
	"os"
	"testing"
)

func TestTerminal(t *testing.T) {
	// test for non-panic
	IsTerminal(os.Stdout.Fd())
}

func TestCygwinPipeName(t *testing.T) {
	if IsCygwinTerminal(os.Stdout.Fd()) {
		t.Fatal("should be false always")
	}
}
