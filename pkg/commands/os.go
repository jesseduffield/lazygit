package commands

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mgutz/str"
	"github.com/sirupsen/logrus"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

// Platform stores the os state
type Platform struct {
	os                   string
	shell                string
	shellArg             string
	escapedQuote         string
	openCommand          string
	openLinkCommand      string
	fallbackEscapedQuote string
}

// OSCommand holds all the os commands
type OSCommand struct {
	Log                *logrus.Entry
	Platform           *Platform
	Config             config.AppConfigurer
	command            func(string, ...string) *exec.Cmd
	getGlobalGitConfig func(string) (string, error)
	getenv             func(string) string
}

// NewOSCommand os command runner
func NewOSCommand(log *logrus.Entry, config config.AppConfigurer) *OSCommand {
	return &OSCommand{
		Log:                log,
		Platform:           getPlatform(),
		Config:             config,
		command:            exec.Command,
		getGlobalGitConfig: gitconfig.Global,
		getenv:             os.Getenv,
	}
}

// RunCommandWithOutput wrapper around commands returning their output and error
func (c *OSCommand) RunCommandWithOutput(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunCommand")
	splitCmd := str.ToArgv(command)
	c.Log.Info(splitCmd)
	return sanitisedCommandOutput(
		c.command(splitCmd[0], splitCmd[1:]...).CombinedOutput(),
	)
}

// RunCommand runs a command and just returns the error
func (c *OSCommand) RunCommand(command string) error {
	_, err := c.RunCommandWithOutput(command)
	return err
}

// FileType tells us if the file is a file, directory or other
func (c *OSCommand) FileType(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "other"
	}
	if fileInfo.IsDir() {
		return "directory"
	}
	return "file"
}

// RunDirectCommand wrapper around direct commands
func (c *OSCommand) RunDirectCommand(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunDirectCommand")

	return sanitisedCommandOutput(
		c.command(c.Platform.shell, c.Platform.shellArg, command).
			CombinedOutput(),
	)
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if err != nil {
		// errors like 'exit status 1' are not very useful so we'll create an error
		// from the combined output
		if outputString == "" {
			return "", err
		}
		return outputString, errors.New(outputString)
	}
	return outputString, nil
}

// OpenFile opens a file with the given
func (c *OSCommand) OpenFile(filename string) error {
	commandTemplate := c.Config.GetUserConfig().GetString("os.openCommand")
	templateValues := map[string]string{
		"filename": c.Quote(filename),
	}

	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	err := c.RunCommand(command)
	return err
}

// OpenLink opens a file with the given
func (c *OSCommand) OpenLink(link string) error {
	commandTemplate := c.Config.GetUserConfig().GetString("os.openLinkCommand")
	templateValues := map[string]string{
		"link": c.Quote(link),
	}

	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	err := c.RunCommand(command)
	return err
}

// EditFile opens a file in a subprocess using whatever editor is available,
// falling back to core.editor, VISUAL, EDITOR, then vi
func (c *OSCommand) EditFile(filename string) (*exec.Cmd, error) {
	editor, _ := c.getGlobalGitConfig("core.editor")

	if editor == "" {
		editor = c.getenv("VISUAL")
	}
	if editor == "" {
		editor = c.getenv("EDITOR")
	}
	if editor == "" {
		if err := c.RunCommand("which vi"); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return nil, errors.New("No editor defined in $VISUAL, $EDITOR, or git config")
	}

	return c.PrepareSubProcess(editor, filename), nil
}

// PrepareSubProcess iniPrepareSubProcessrocess then tells the Gui to switch to it
func (c *OSCommand) PrepareSubProcess(cmdName string, commandArgs ...string) *exec.Cmd {
	return c.command(cmdName, commandArgs...)
}

// Quote wraps a message in platform-specific quotation marks
func (c *OSCommand) Quote(message string) string {
	message = strings.Replace(message, "`", "\\`", -1)
	escapedQuote := c.Platform.escapedQuote
	if strings.Contains(message, c.Platform.escapedQuote) {
		escapedQuote = c.Platform.fallbackEscapedQuote
	}
	return escapedQuote + message + escapedQuote
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
