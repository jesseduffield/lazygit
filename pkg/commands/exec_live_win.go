// +build windows

package commands

// RunCommandWithOutputLiveWrapper runs a command live but because of windows compatibility this command can't be ran there
func RunCommandWithOutputLiveWrapper(c *OSCommand, command string, output func(string) string) error {
	return c.RunCommand(command)
}
