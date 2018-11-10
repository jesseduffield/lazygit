// +build windows

package commands

// RunCommandWithOutputLiveWrapper runs a command live but because of windows compatibility this command can't be ran there
// TODO: Remove this hack and replace it with a propper way to run commands live on windows
func RunCommandWithOutputLiveWrapper(c *OSCommand, command string, output func(string) string) (errorMessage string, codeError error) {
	cmdOputput := c.RunCommand(command)
	return cmdOputput.Error(), cmdOputput
}
