// +build windows

package oscommands

// runAndParseWords runs a command live but because of windows compatibility this command can't be ran there
// TODO: Remove this hack and replace it with a proper way to run commands live on windows
func runAndParseWords(c *OS, cmdObj ICmdObj, output func(string) string) error {
	return c.Run(cmdObj)
}
