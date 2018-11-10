// +build !windows

package commands

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/kr/pty"
)

// RunCommandWithOutputLiveWrapper runs a command and return every word that gets written in stdout
// Output is a function that executes by every word that gets read by bufio
// As return of output you need to give a string that will be written to stdin
// NOTE: If the return data is empty it won't written anything to stdin
// NOTE: You don't have to include a enter in the return data this function will do that for you
func RunCommandWithOutputLiveWrapper(c *OSCommand, command string, output func(string) string) (errorMessage string, codeError error) {
	cmdOutput := []string{}

	splitCmd := ToArgv(command)
	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=en_US.utf8", "LC_ALL=en_US.UTF-8")

	tty, err := pty.Start(cmd)

	if err != nil {
		return errorMessage, err
	}

	stopAsking := make(chan struct{})

	var waitForBufio sync.WaitGroup
	waitForBufio.Add(1)

	defer func() {
		_ = tty.Close()
	}()

	go func() {
		scanner := bufio.NewScanner(tty)
		scanner.Split(scanWordsWithNewLines)
	loop:
		for scanner.Scan() {
			select {
			case <-stopAsking:
				break loop
			default:
				toOutput := strings.Trim(scanner.Text(), " ")
				cmdOutput = append(cmdOutput, toOutput)
				toWrite := output(toOutput)
				if len(toWrite) > 0 {
					_, _ = tty.Write([]byte(toWrite + "\n"))
				}
			}
		}
		waitForBufio.Done()
	}()

	if err = cmd.Wait(); err != nil {
		stopAsking <- struct{}{}
		waitForBufio.Wait()
		return strings.Join(cmdOutput, " "), err
	}

	return errorMessage, nil
}

// scanWordsWithNewLines is a copy of bufio.ScanWords but this also captures new lines
// For specific comments about this function take a look at: bufio.ScanWords
func scanWordsWithNewLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			return i + width, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return start, nil, nil
}

// isSpace is also copied from the bufio package and has been modified to also captures new lines
// For specific comments about this function take a look at: bufio.isSpace
func isSpace(r rune) bool {
	if r <= '\u00FF' {
		switch r {
		case ' ', '\t', '\v', '\f':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}
