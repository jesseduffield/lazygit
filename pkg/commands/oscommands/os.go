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
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Platform stores the os state
type Platform struct {
	OS              string
	Shell           string
	ShellArg        string
	OpenCommand     string
	OpenLinkCommand string
}

type ICommander interface {
	Run(ICmdObj) error
	RunWithOutput(ICmdObj) (string, error)
}

type RealCommander struct {
}

func (self *RealCommander) Run(cmdObj ICmdObj) error {
	return cmdObj.GetCmd().Run()
}

// OSCommand holds all the os commands
type OSCommand struct {
	*common.Common
	Platform *Platform
	Command  func(string, ...string) *exec.Cmd
	Getenv   func(string) string

	// callback to run before running a command, i.e. for the purposes of logging
	onRunCommand func(CmdLogEntry)

	// something like 'Staging File': allows us to group cmd logs under a single title
	CmdLogSpan string

	removeFile func(string) error

	Cmd *CmdObjBuilder
}

// TODO: make these fields private
type CmdLogEntry struct {
	// e.g. 'git commit -m "haha"'
	cmdStr string
	// Span is something like 'Staging File'. Multiple commands can be grouped under the same
	// span
	span string

	// sometimes our command is direct like 'git commit', and sometimes it's a
	// command to remove a file but through Go's standard library rather than the
	// command line
	commandLine bool
}

func (e CmdLogEntry) GetCmdStr() string {
	return e.cmdStr
}

func (e CmdLogEntry) GetSpan() string {
	return e.span
}

func (e CmdLogEntry) GetCommandLine() bool {
	return e.commandLine
}

func NewCmdLogEntry(cmdStr string, span string, commandLine bool) CmdLogEntry {
	return CmdLogEntry{cmdStr: cmdStr, span: span, commandLine: commandLine}
}

// NewOSCommand os command runner
func NewOSCommand(common *common.Common) *OSCommand {
	command := secureexec.Command
	platform := getPlatform()

	c := &OSCommand{
		Common:     common,
		Platform:   platform,
		Command:    command,
		Getenv:     os.Getenv,
		removeFile: os.RemoveAll,
	}

	runner := &RealRunner{log: common.Log, logCmdObj: c.LogCmdObj}
	c.Cmd = &CmdObjBuilder{runner: runner, command: command, logCmdObj: c.LogCmdObj, platform: platform}

	return c
}

func (c *OSCommand) WithSpan(span string) *OSCommand {
	// sometimes .WithSpan(span) will be called where span actually is empty, in
	// which case we don't need to log anything so we can just return early here
	// with the original struct
	if span == "" {
		return c
	}

	newOSCommand := &OSCommand{}
	*newOSCommand = *c
	newOSCommand.CmdLogSpan = span
	return newOSCommand
}

func (c *OSCommand) LogCmdObj(cmdObj ICmdObj) {
	c.LogCommand(cmdObj.ToString(), true)
}

func (c *OSCommand) LogCommand(cmdStr string, commandLine bool) {
	c.Log.WithField("command", cmdStr).Info("RunCommand")

	if c.onRunCommand != nil && c.CmdLogSpan != "" {
		c.onRunCommand(NewCmdLogEntry(cmdStr, c.CmdLogSpan, commandLine))
	}
}

func (c *OSCommand) SetOnRunCommand(f func(CmdLogEntry)) {
	c.onRunCommand = f
}

// SetCommand sets the command function used by the struct.
// To be used for testing only
func (c *OSCommand) SetCommand(cmd func(string, ...string) *exec.Cmd) {
	c.Command = cmd
}

// To be used for testing only
func (c *OSCommand) SetRemoveFile(f func(string) error) {
	c.removeFile = f
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

// OpenFile opens a file with the given
func (c *OSCommand) OpenFile(filename string) error {
	commandTemplate := c.UserConfig.OS.OpenCommand
	templateValues := map[string]string{
		"filename": c.Quote(filename),
	}
	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	return c.Cmd.NewShell(command).Run()
}

// OpenLink opens a file with the given
func (c *OSCommand) OpenLink(link string) error {
	c.LogCommand(fmt.Sprintf("Opening link '%s'", link), false)
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

func (self *CmdObjBuilder) Quote(message string) string {
	var quote string
	if self.platform.OS == "windows" {
		quote = `\"`
		message = strings.NewReplacer(
			`"`, `"'"'"`,
			`\"`, `\\"`,
		).Replace(message)
	} else {
		quote = `"`
		message = strings.NewReplacer(
			`\`, `\\`,
			`"`, `\"`,
			`$`, `\$`,
			"`", "\\`",
		).Replace(message)
	}
	return quote + message + quote
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

// GetLazygitPath returns the path of the currently executed file
func (c *OSCommand) GetLazygitPath() string {
	ex, err := os.Executable() // get the executable path for git to use
	if err != nil {
		ex = os.Args[0] // fallback to the first call argument if needed
	}
	return `"` + filepath.ToSlash(ex) + `"`
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

func GetTempDir() string {
	return filepath.Join(os.TempDir(), "lazygit")
}
