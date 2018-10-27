// +build !windows

package commands

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"regexp"

	"github.com/kr/pty"
	"github.com/mgutz/str"
)

// RunCommandWithOutputLiveWrapper runs a command and return every word that gets written in stdout
// Output is a function that executes by every word that gets read by bufio
// As return of output you need to give a string that will be written to stdin
// NOTE: If the return data is empty it won't written anything to stdin
// NOTE: You don't have to include a enter in the return data this function will do that for you
func RunCommandWithOutputLiveWrapper(c *OSCommand, command string, output func(string) string) error {
	splitCmd := str.ToArgv(command)
	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=en_US.utf8", "LC_ALL=en_US.UTF-8")

	tty, err := pty.Start(cmd)

	if err != nil {
		return errors.New(err.Error())
	}

	defer func() { _ = tty.Close() }()

	go func() {
		// Regex to cleanup the command output
		// sometimes the output words include unneeded spaces at eatch end of the string
		re := regexp.MustCompile(`(^\s*)|(\s*$)`)

		scanner := bufio.NewScanner(tty)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			toWrite := output(re.ReplaceAllString(scanner.Text(), ""))
			if len(toWrite) > 0 {
				_, _ = tty.Write([]byte(toWrite + "\n"))
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
