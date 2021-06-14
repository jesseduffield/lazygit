package oscommands

import (
	"bufio"
	"os"
	"os/exec"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/mgutz/str"
)

func (c *OSCommand) RunCommandWithOutput(cmdObj ICmdObj) (string, error) {
	c.LogCmd(cmdObj)
	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		c.Log.WithField("command", cmdObj.ToString()).Error(output)
	}
	return output, err
}

// RunExecutable runs an executable file and returns an error if there was one
func (c *OSCommand) RunExecutable(cmd ICmdObj) error {
	_, err := c.RunCommandWithOutput(cmd)

	return err
}

// ExecutableFromString takes a string like `git status` and returns an executable command for it
func (c *OSCommand) ExecutableFromString(commandStr string) *exec.Cmd {
	splitCmd := str.ToArgv(commandStr)
	cmd := c.Command(splitCmd[0], splitCmd[1:]...)
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	return cmd
}

func (c *OSCommand) RunCommandAndParseOutput(cmdObj ICmdObj, output func(string) string) error {
	return runCommandAndParseOutput(c, cmdObj, output)
}

func RunLineOutputCmd(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
	cmd := cmdObj.GetCmd()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutPipe)
	scanner.Split(bufio.ScanLines)
	if err := cmd.Start(); err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		stop, err := onLine(line)
		if err != nil {
			return err
		}
		if stop {
			_ = cmd.Process.Kill()
			break
		}
	}

	_ = cmd.Wait()

	return nil
}
