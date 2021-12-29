package oscommands

import (
	"bufio"
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
	"github.com/mgutz/str"
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

	IRunner
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
	c := &OSCommand{
		Common:     common,
		Platform:   getPlatform(),
		Command:    secureexec.Command,
		Getenv:     os.Getenv,
		removeFile: os.RemoveAll,
	}

	c.IRunner = &RealRunner{c: c}
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
	err := c.Run(c.NewShellCmdObjFromString(command))
	return err
}

// OpenLink opens a file with the given
func (c *OSCommand) OpenLink(link string) error {
	c.LogCommand(fmt.Sprintf("Opening link '%s'", link), false)
	commandTemplate := c.UserConfig.OS.OpenLinkCommand
	templateValues := map[string]string{
		"link": c.Quote(link),
	}

	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	err := c.Run(c.NewShellCmdObjFromString(command))
	return err
}

// Quote wraps a message in platform-specific quotation marks
func (c *OSCommand) Quote(message string) string {
	var quote string
	if c.Platform.OS == "windows" {
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
		cmds[i] = c.NewCmdObj(str).GetCmd()
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

// builders

func (c *OSCommand) NewCmdObj(cmdStr string) ICmdObj {
	args := str.ToArgv(cmdStr)
	cmd := c.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	return &CmdObj{
		cmdStr: cmdStr,
		cmd:    cmd,
	}
}

func (c *OSCommand) NewCmdObjFromArgs(args []string) ICmdObj {
	cmd := c.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	return &CmdObj{
		cmdStr: strings.Join(args, " "),
		cmd:    cmd,
	}
}

// NewShellCmdObjFromString takes a string like `git commit` and returns an executable shell command for it
func (c *OSCommand) NewShellCmdObjFromString(commandStr string) ICmdObj {
	quotedCommand := ""
	// Windows does not seem to like quotes around the command
	if c.Platform.OS == "windows" {
		quotedCommand = strings.NewReplacer(
			"^", "^^",
			"&", "^&",
			"|", "^|",
			"<", "^<",
			">", "^>",
			"%", "^%",
		).Replace(commandStr)
	} else {
		quotedCommand = c.Quote(commandStr)
	}

	shellCommand := fmt.Sprintf("%s %s %s", c.Platform.Shell, c.Platform.ShellArg, quotedCommand)
	return c.NewCmdObj(shellCommand)
}

// TODO: pick one of NewShellCmdObjFromString2 and ShellCommandFromString to use. I'm not sure
// which one actually is better, but I suspect it's NewShellCmdObjFromString2
func (c *OSCommand) NewShellCmdObjFromString2(command string) ICmdObj {
	return c.NewCmdObjFromArgs([]string{c.Platform.Shell, c.Platform.ShellArg, command})
}

// runners

type IRunner interface {
	Run(cmdObj ICmdObj) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	RunLineOutputCmd(cmdObj ICmdObj, onLine func(line string) (bool, error)) error
}

type RunExpectation func(ICmdObj) (string, error)

type RealRunner struct {
	c *OSCommand
}

func (self *RealRunner) Run(cmdObj ICmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *RealRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	self.c.LogCmdObj(cmdObj)
	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		self.c.Log.WithField("command", cmdObj.ToString()).Error(output)
	}
	return output, err
}

func (self *RealRunner) RunLineOutputCmd(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
	cmd := cmdObj.GetCmd()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutPipe)
	scanner.Split(bufio.ScanLines)
	if err := cmd.Start(); err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		stop, err := onLine(line)
		if err != nil {
			return err
		}
		if stop {
			_ = cmd.Process.Kill()
			break
		}
	}

	_ = cmd.Wait()

	return nil
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if err != nil {
		// errors like 'exit status 1' are not very useful so we'll create an error
		// from the combined output
		if outputString == "" {
			return "", utils.WrapError(err)
		}
		return outputString, errors.New(outputString)
	}
	return outputString, nil
}
