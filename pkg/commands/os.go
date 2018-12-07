package commands

import (
	"errors"
	"io/ioutil"
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
	Log      *logrus.Entry
	Platform *Platform
	Config   config.AppConfigurer
}

// regenerate the mock for this interface with
// mockgen -source=os.go -destination=mock_os.go  -package=commands MockOSCommand
// Command
type Command interface {
	RunCommandWithOutput(command string) (string, error)
	RunCommand(command string) error
	FileType(path string) string
	RunDirectCommand(command string) (string, error)
	OpenFile(filename string) error
	OpenLink(link string) error
	EditFile(filename string) (*exec.Cmd, error)
	PrepareSubProcess(cmdName string, commandArgs ...string) *exec.Cmd
	Quote(message string) string
	Unquote(message string) string
	AppendLineToFile(filename, line string) error
	CreateTempFile(filename, content string) (string, error)
	RemoveFile(filename string) error
	GetPlatform() *Platform
	SetGetGlobalGitConfig(func(string) (string, error))
	GetEnv(string) string
	Command(string, ...string) *exec.Cmd
}

// NewOSCommand os command runner
func NewOSCommand(log *logrus.Entry, config config.AppConfigurer) *OSCommand {
	osCommand := &OSCommand{
		Log:      log,
		Platform: getPlatform(),
		Config:   config,
	}

	return osCommand
}

// RunCommandWithOutput wrapper around commands returning their output and error
func (c *OSCommand) RunCommandWithOutput(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunCommand")
	splitCmd := str.ToArgv(command)
	c.Log.Info(splitCmd)
	return sanitisedCommandOutput(
		c.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput(),
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
		c.Command(c.Platform.shell, c.Platform.shellArg, command).
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
	editor, _ := c.GetGlobalGitConfig("core.editor")

	if editor == "" {
		editor = c.GetEnv("VISUAL")
	}
	if editor == "" {
		editor = c.GetEnv("EDITOR")
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
	return c.Command(cmdName, commandArgs...)
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

// CreateTempFile writes a string to a new temp file and returns the file's name
func (c *OSCommand) CreateTempFile(filename, content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		c.Log.Error(err)
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		c.Log.Error(err)
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		c.Log.Error(err)
		return "", err
	}

	return tmpfile.Name(), nil
}

// RemoveFile removes a file at the specified path
func (c *OSCommand) RemoveFile(filename string) error {
	return os.Remove(filename)
}

// GetPlatform returns the computer's platform
func (c *OSCommand) GetPlatform() *Platform {
	return c.Platform
}

func (c *OSCommand) GetGlobalGitConfig(s string) (string, error) {
	return gitconfig.Global(s)
}

func (c *OSCommand) GetEnv(s string) string {
	return os.Getenv(s)
}

func (c *OSCommand) Command(s string, s2 ...string) *exec.Cmd {
	return exec.Command(s, s2...)
}
