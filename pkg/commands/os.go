package commands

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jesseduffield/gocui"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

var (
	// ErrNoOpenCommand : When we don't know which command to use to open a file
	ErrNoOpenCommand = errors.New("Unsure what command to use to open this file")
)

// Platform stores the os state
type platform struct {
	os           string
	shell        string
	shellArg     string
	escapedQuote string
}

// OSCommand holds all the os commands
type OSCommand struct {
	Log      *logrus.Logger
	Platform platform
}

// NewOSCommand os command runner
func NewOSCommand(log *logrus.Logger) (*OSCommand, error) {
	osCommand := &OSCommand{
		Log:      log,
		Platform: getPlatform(),
	}
	return osCommand, nil
}

// RunCommand wrapper around commands
func (c *OSCommand) RunCommand(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunCommand")
	splitCmd := strings.Split(command, " ")
	cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput()
	return sanitisedCommandOutput(cmdOut, err)
}

// RunDirectCommand wrapper around direct commands
func (c *OSCommand) RunDirectCommand(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunDirectCommand")

	cmdOut, err := exec.
		Command(c.Platform.shell, c.Platform.shellArg, command).
		CombinedOutput()
	return sanitisedCommandOutput(cmdOut, err)
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if outputString == "" && err != nil {
		return err.Error(), err
	}
	return outputString, err
}

func getPlatform() platform {
	switch runtime.GOOS {
	case "windows":
		return platform{
			os:           "windows",
			shell:        "cmd",
			shellArg:     "/c",
			escapedQuote: "\\\"",
		}
	default:
		return platform{
			os:           runtime.GOOS,
			shell:        "bash",
			shellArg:     "-c",
			escapedQuote: "\"",
		}
	}
}

func (c *OSCommand) getOpenCommand() (string, string, error) {
	//NextStep open equivalents: xdg-open (linux), cygstart (cygwin), open (OSX)
	trailMap := map[string]string{
		"xdg-open": " &>/dev/null &",
		"cygstart": "",
		"open":     "",
	}
	for name, trail := range trailMap {
		if out, _ := c.runCommand("which " + name); out != "exit status 1" {
			return name, trail, nil
		}
	}
	return "", "", ErrNoOpenCommand
}

// VsCodeOpenFile opens the file in code, with the -r flag to open in the
// current window
func (c *OSCommand) VsCodeOpenFile(g *gocui.Gui, filename string) (string, error) {
	return c.runCommand("code -r " + filename)
}

// SublimeOpenFile opens the filein sublime
// may be deprecated in the future
func (c *OSCommand) SublimeOpenFile(g *gocui.Gui, filename string) (string, error) {
	return c.runCommand("subl " + filename)
}

// OpenFile opens a file with the given
func (c *OSCommand) OpenFile(g *gocui.Gui, filename string) (string, error) {
	cmdName, cmdTrail, err := getOpenCommand()
	if err != nil {
		return "", err
	}
	return c.runCommand(cmdName + " " + filename + cmdTrail)
}

// EditFile opens a file in a subprocess using whatever editor is available,
// falling back to core.editor, VISUAL, EDITOR, then vi
func (c *OSCommand) editFile(g *gocui.Gui, filename string) (string, error) {
	editor, _ := gitconfig.Global("core.editor")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		if _, err := c.OSCommand.runCommand("which vi"); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return "", createErrorPanel(g, "No editor defined in $VISUAL, $EDITOR, or git config.")
	}
	runSubProcess(g, editor, filename)
	return "", nil
}
