// +build windows

package sequences

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestStdoutSequencesOn(t *testing.T) {
	err := EnableVirtualTerminalProcessing(syscall.Stdout, true)
	if err != nil {
		t.Fatalf("Failed to enable VTP: %v", err)
	}
	defer EnableVirtualTerminalProcessing(syscall.Stdout, false)

	fmt.Fprintf(os.Stdout, "\x1b[34mHello \x1b[35mWorld\x1b[0m!\n")
}

func TestStdoutSequencesOff(t *testing.T) {
	err := EnableVirtualTerminalProcessing(syscall.Stdout, false)
	if err != nil {
		t.Fatalf("Failed to enable VTP: %v", err)
	}

	fmt.Fprintf(os.Stdout, "\x1b[34mHello \x1b[35mWorld\x1b[0m!\n")
}

func TestStderrSequencesOn(t *testing.T) {
	err := EnableVirtualTerminalProcessing(syscall.Stderr, true)
	if err != nil {
		t.Fatalf("Failed to enable VTP: %v", err)
	}
	defer EnableVirtualTerminalProcessing(syscall.Stderr, false)

	fmt.Fprintf(os.Stderr, "\x1b[34mHello \x1b[35mWorld\x1b[0m!\n")
}

func TestStderrSequencesOff(t *testing.T) {
	err := EnableVirtualTerminalProcessing(syscall.Stderr, false)
	if err != nil {
		t.Fatalf("Failed to enable VTP: %v", err)
	}

	fmt.Fprintf(os.Stderr, "\x1b[34mHello \x1b[35mWorld\x1b[0m!\n")
}
