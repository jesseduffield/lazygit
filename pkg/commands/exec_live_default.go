// +build !windows

package commands

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/go-errors/errors"

	"github.com/mgutz/str"
)

// RunCommandWithOutputLiveWrapper runs a command and return every word that gets written in stdout
// Output is a function that executes by every word that gets read by bufio
// As return of output you need to give a string that will be written to stdin
// NOTE: If the return data is empty it won't written anything to stdin
func RunCommandWithOutputLiveWrapper(c *OSCommand, command string, output func(string) string) error {
	splitCmd := str.ToArgv(command)
	cmd := c.command(splitCmd[0], splitCmd[1:]...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanRunes)
		for scanner.Scan() {
			toOutput := scanner.Text()
			credential := output(toOutput)
			if credential != "" {
				if _, err := io.WriteString(stdin, credential); err != nil {
					c.Log.Error(err)
				}
			}
		}
	}()

	if err = cmd.Wait(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
