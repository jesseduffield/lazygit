// +build windows

package oscommands

// runCommandAndParseOutput runs a command live but because of windows compatibility this command can't be ran there
// TODO: Remove this hack and replace it with a proper way to run commands live on windows
func runCommandAndParseOutput(c *OSCommand, cmdObj *CmdObj, output func(string) string) error {
	return c.RunCommand(cmdObj.ToString())
}
