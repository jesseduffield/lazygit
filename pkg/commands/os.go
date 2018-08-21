package commands

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/mgutz/str"

	"github.com/sirupsen/logrus"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

// Platform stores the os state
type Platform struct {
	os           string
	shell        string
	shellArg     string
	escapedQuote string
}

// OSCommand holds all the os commands
type OSCommand struct {
	Log      *logrus.Logger
	Platform *Platform
}

// NewOSCommand os command runner
func NewOSCommand(log *logrus.Logger) *OSCommand {
	return &OSCommand{
		Log:      log,
		Platform: getPlatform(),
	}
}

// RunCommandWithOutput wrapper around commands returning their output and error
func (c *OSCommand) RunCommandWithOutput(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunCommand")
	splitCmd := str.ToArgv(command)
	c.Log.Info(splitCmd)

	return sanitisedCommandOutput(
		exec.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput(),
	)
}

// RunCommand runs a command and just returns the error
func (c *OSCommand) RunCommand(command string) error {
	_, err := c.RunCommandWithOutput(command)
	return err
}

// RunDirectCommand wrapper around direct commands
func (c *OSCommand) RunDirectCommand(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunDirectCommand")
	args := str.ToArgv(c.Platform.shellArg + " " + command)
	c.Log.Info(spew.Sdump(args))

	return sanitisedCommandOutput(
		exec.
			Command(c.Platform.shell, args...).
			CombinedOutput(),
	)
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if err != nil {
		// errors like 'exit status 1' are not very useful so we'll create an error
		// from the combined output
		return outputString, errors.New(outputString)
	}
	return outputString, nil
}

// GetOpenCommand get open command
func (c *OSCommand) GetOpenCommand() (string, string, error) {
	//NextStep open equivalents: xdg-open (linux), cygstart (cygwin), open (OSX)
	trailMap := map[string]string{
		"xdg-open": " &>/dev/null &",
		"cygstart": "",
		"open":     "",
	}
	for name, trail := range trailMap {
		if err := c.RunCommand("which " + name); err == nil {
			return name, trail, nil
		}
	}
	return "", "", errors.New("Unsure what command to use to open this file")
}

// VsCodeOpenFile opens the file in code, with the -r flag to open in the
// current window
// each of these open files needs to have the same function signature because
// they're being passed as arguments into another function,
// but only editFile actually returns a *exec.Cmd
func (c *OSCommand) VsCodeOpenFile(filename string) (*exec.Cmd, error) {
	return nil, c.RunCommand("code -r " + filename)
}

// SublimeOpenFile opens the filein sublime
// may be deprecated in the future
func (c *OSCommand) SublimeOpenFile(filename string) (*exec.Cmd, error) {
	return nil, c.RunCommand("subl " + filename)
}

// OpenFile opens a file with the given
func (c *OSCommand) OpenFile(filename string) (*exec.Cmd, error) {
	cmdName, cmdTrail, err := c.GetOpenCommand()
	if err != nil {
		return nil, err
	}
	err = c.RunCommand(cmdName + " " + c.Quote(filename) + cmdTrail) // TODO: test on linux
	return nil, err
}

// EditFile opens a file in a subprocess using whatever editor is available,
// falling back to core.editor, VISUAL, EDITOR, then vi
func (c *OSCommand) EditFile(filename string) (*exec.Cmd, error) {
	editor, _ := gitconfig.Global("core.editor")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		if err := c.RunCommand("which vi"); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return nil, errors.New("No editor defined in $VISUAL, $EDITOR, or git config")
	}
	return c.PrepareSubProcess(editor, filename)
}

// PrepareSubProcess iniPrepareSubProcessrocess then tells the Gui to switch to it
func (c *OSCommand) PrepareSubProcess(cmdName string, commandArgs ...string) (*exec.Cmd, error) {
	subprocess := exec.Command(cmdName, commandArgs...)
	return subprocess, nil
}

// Quote wraps a message in platform-specific quotation marks
func (c *OSCommand) Quote(message string) string {
	message = strings.Replace(message, "`", "\\`", -1)
	return c.Platform.escapedQuote + message + c.Platform.escapedQuote
}

// Unquote removes wrapping quotations marks if they are present
// this is needed for removing quotes from staged filenames with spaces
func (c *OSCommand) Unquote(message string) string {
	return strings.Replace(message, `"`, "", -1)
}

// AppendLineToFile adds a new line in file
func (c *OSCommand) AppendLineToFile(filename, line string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + line)
	return err
}
