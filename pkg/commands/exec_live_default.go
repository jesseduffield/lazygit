// +build !windows

package commands

import (
	"bufio"
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/go-errors/errors"

	"github.com/jesseduffield/pty"
)

// RunCommandWithOutputLiveWrapper runs a command and return every word that gets written in stdout
// Output is a function that executes by every word that gets read by bufio
// As return of output you need to give a string that will be written to stdin
// NOTE: If the return data is empty it won't written anything to stdin
func RunCommandWithOutputLiveWrapper(c *OSCommand, command string, output func(string) string) error {
	cmd := c.ExecutableFromString(command)
	cmd.Env = append(cmd.Env, "LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	ptmx, err := pty.Start(cmd)

	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(ptmx)
		scanner.Split(scanWordsWithNewLines)
		for scanner.Scan() {
			toOutput := strings.Trim(scanner.Text(), " ")
			_, _ = ptmx.WriteString(output(toOutput))
		}
	}()

	err = cmd.Wait()
	ptmx.Close()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
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
