package oscommands

import (
	"bufio"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (c *OSCommand) RunWithOutput(cmdObj ICmdObj) (string, error) {
	c.LogCmd(cmdObj)
	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		c.Log.WithField("command", cmdObj.ToString()).Error(output)
	}
	return output, err
}

// Run runs an executable file and returns an error if there was one
func (c *OSCommand) Run(cmd ICmdObj) error {
	_, err := c.RunWithOutput(cmd)

	return err
}

func (c *OSCommand) RunAndParseWords(cmdObj ICmdObj, output func(string) string) error {
	return runAndParseWords(c, cmdObj, output)
}

func RunAndParseLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
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
