package oscommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-errors/errors"

	"github.com/atotto/clipboard"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// OSCommand holds all the os commands
type OSCommand struct {
	*common.Common
	Platform *Platform
	GetenvFn func(string) string

	// callback to run before running a command, i.e. for the purposes of logging.
	// the string argument is the command string e.g. 'git add .' and the bool is
	// whether we're dealing with a command line command or something more general
	// like 'Opening PR URL', or something handled by Go's standard library.
	logCommandFn func(string, bool)

	removeFile func(string) error

	Cmd *CmdObjBuilder
}

// Platform stores the os state
type Platform struct {
	OS              string
	Shell           string
	ShellArg        string
	OpenCommand     string
	OpenLinkCommand string
}

// NewOSCommand os command runner
func NewOSCommand(common *common.Common, platform *Platform, guiIO *guiIO) *OSCommand {
	c := &OSCommand{
		Common:     common,
		Platform:   platform,
		GetenvFn:   os.Getenv,
		removeFile: os.RemoveAll,
	}

	runner := &cmdObjRunner{log: common.Log, guiIO: guiIO}
	c.Cmd = &CmdObjBuilder{runner: runner, platform: platform}

	return c
}

func (c *OSCommand) LogCommand(cmdStr string, commandLine bool) {
	c.Log.WithField("command", cmdStr).Info("RunCommand")

	if c.logCommandFn != nil {
		c.logCommandFn(cmdStr, commandLine)
	}
}

func (c *OSCommand) SetLogCommandFn(f func(string, bool)) {
	c.logCommandFn = f
}

// To be used for testing only
func (c *OSCommand) SetRemoveFile(f func(string) error) {
	c.removeFile = f
}

// FileType tells us if the file is a file, directory or other
func FileType(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "other"
	}
	if fileInfo.IsDir() {
		return "directory"
	}
	return "file"
}

func (c *OSCommand) OpenFile(filename string) error {
	commandTemplate := c.UserConfig.OS.OpenCommand
	templateValues := map[string]string{
		"filename": c.Quote(filename),
	}
	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	return c.Cmd.NewShell(command).Run()
}

func (c *OSCommand) OpenLink(link string) error {
	commandTemplate := c.UserConfig.OS.OpenLinkCommand
	templateValues := map[string]string{
		"link": c.Quote(link),
	}

	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	return c.Cmd.NewShell(command).Run()
}

// Quote wraps a message in platform-specific quotation marks
func (c *OSCommand) Quote(message string) string {
	return c.Cmd.Quote(message)
}

// AppendLineToFile adds a new line in file
func (c *OSCommand) AppendLineToFile(filename, line string) error {
	c.LogCommand(fmt.Sprintf("Appending '%s' to file '%s'", line, filename), false)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return utils.WrapError(err)
	}
	defer f.Close()

	_, err = f.WriteString("\n" + line)
	if err != nil {
		return utils.WrapError(err)
	}
	return nil
}

// CreateTempFile writes a string to a new temp file and returns the file's name
func (c *OSCommand) CreateTempFile(filename, content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		c.Log.Error(err)
		return "", utils.WrapError(err)
	}
	c.LogCommand(fmt.Sprintf("Creating temp file '%s'", tmpfile.Name()), false)

	if _, err := tmpfile.WriteString(content); err != nil {
		c.Log.Error(err)
		return "", utils.WrapError(err)
	}
	if err := tmpfile.Close(); err != nil {
		c.Log.Error(err)
		return "", utils.WrapError(err)
	}

	return tmpfile.Name(), nil
}

// CreateFileWithContent creates a file with the given content
func (c *OSCommand) CreateFileWithContent(path string, content string) error {
	c.LogCommand(fmt.Sprintf("Creating file '%s'", path), false)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		c.Log.Error(err)
		return err
	}

	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		c.Log.Error(err)
		return utils.WrapError(err)
	}

	return nil
}

// Remove removes a file or directory at the specified path
func (c *OSCommand) Remove(filename string) error {
	c.LogCommand(fmt.Sprintf("Removing '%s'", filename), false)
	err := os.RemoveAll(filename)
	return utils.WrapError(err)
}

// FileExists checks whether a file exists at the specified path
func (c *OSCommand) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// PipeCommands runs a heap of commands and pipes their inputs/outputs together like A | B | C
func (c *OSCommand) PipeCommands(commandStrings ...string) error {
	cmds := make([]*exec.Cmd, len(commandStrings))
	logCmdStr := ""
	for i, str := range commandStrings {
		if i > 0 {
			logCmdStr += " | "
		}
		logCmdStr += str
		cmds[i] = c.Cmd.New(str).GetCmd()
	}
	c.LogCommand(logCmdStr, true)

	for i := 0; i < len(cmds)-1; i++ {
		stdout, err := cmds[i].StdoutPipe()
		if err != nil {
			return err
		}

		cmds[i+1].Stdin = stdout
	}

	// keeping this here in case I adapt this code for some other purpose in the future
	// cmds[len(cmds)-1].Stdout = os.Stdout

	finalErrors := []string{}

	wg := sync.WaitGroup{}
	wg.Add(len(cmds))

	for _, cmd := range cmds {
		currentCmd := cmd
		go utils.Safe(func() {
			stderr, err := currentCmd.StderrPipe()
			if err != nil {
				c.Log.Error(err)
			}

			if err := currentCmd.Start(); err != nil {
				c.Log.Error(err)
			}

			if b, err := ioutil.ReadAll(stderr); err == nil {
				if len(b) > 0 {
					finalErrors = append(finalErrors, string(b))
				}
			}

			if err := currentCmd.Wait(); err != nil {
				c.Log.Error(err)
			}

			wg.Done()
		})
	}

	wg.Wait()

	if len(finalErrors) > 0 {
		return errors.New(strings.Join(finalErrors, "\n"))
	}
	return nil
}

func Kill(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		// somebody got to it before we were able to, poor bastard
		return nil
	}
	return cmd.Process.Kill()
}

func (c *OSCommand) CopyToClipboard(str string) error {
	escaped := strings.Replace(str, "\n", "\\n", -1)
	truncated := utils.TruncateWithEllipsis(escaped, 40)
	c.LogCommand(fmt.Sprintf("Copying '%s' to clipboard", truncated), false)
	return clipboard.WriteAll(str)
}

func (c *OSCommand) RemoveFile(path string) error {
	c.LogCommand(fmt.Sprintf("Deleting path '%s'", path), false)

	return c.removeFile(path)
}

func (c *OSCommand) Getenv(key string) string {
	return c.GetenvFn(key)
}

func GetTempDir() string {
	return filepath.Join(os.TempDir(), "lazygit")
}

// GetLazygitPath returns the path of the currently executed file
func GetLazygitPath() string {
	ex, err := os.Executable() // get the executable path for git to use
	if err != nil {
		ex = os.Args[0] // fallback to the first call argument if needed
	}
	return `"` + filepath.ToSlash(ex) + `"`
}
