package oscommands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-errors/errors"

	"github.com/atotto/clipboard"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mgutz/str"
	"github.com/sirupsen/logrus"
)

// Platform stores the os state
type Platform struct {
	OS              string
	CatCmd          []string
	Shell           string
	ShellArg        string
	EscapedQuote    string
	OpenCommand     string
	OpenLinkCommand string
}

// OSCommand holds all the os commands
type OSCommand struct {
	Log              *logrus.Entry
	Platform         *Platform
	Config           config.AppConfigurer
	Command          func(string, ...string) *exec.Cmd
	BeforeExecuteCmd func(*exec.Cmd)
	Getenv           func(string) string

	// callback to run before running a command, i.e. for the purposes of logging
	onRunCommand func(CmdLogEntry)

	// something like 'Staging File': allows us to group cmd logs under a single title
	CmdLogSpan string

	removeFile func(string) error
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
func NewOSCommand(log *logrus.Entry, config config.AppConfigurer) *OSCommand {
	return &OSCommand{
		Log:              log,
		Platform:         getPlatform(),
		Config:           config,
		Command:          secureexec.Command,
		BeforeExecuteCmd: func(*exec.Cmd) {},
		Getenv:           os.Getenv,
		removeFile:       os.RemoveAll,
	}
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

func (c *OSCommand) LogExecCmd(cmd *exec.Cmd) {
	c.LogCommand(strings.Join(cmd.Args, " "), true)
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

func (c *OSCommand) SetBeforeExecuteCmd(cmd func(*exec.Cmd)) {
	c.BeforeExecuteCmd = cmd
}

type RunCommandOptions struct {
	EnvVars []string
}

func (c *OSCommand) RunCommandWithOutputWithOptions(command string, options RunCommandOptions) (string, error) {
	c.LogCommand(command, true)
	cmd := c.ExecutableFromString(command)

	cmd.Env = append(cmd.Env, "GIT_TERMINAL_PROMPT=0") // prevents git from prompting us for input which would freeze the program
	cmd.Env = append(cmd.Env, options.EnvVars...)

	return sanitisedCommandOutput(cmd.CombinedOutput())
}

func (c *OSCommand) RunCommandWithOptions(command string, options RunCommandOptions) error {
	_, err := c.RunCommandWithOutputWithOptions(command, options)
	return err
}

// RunCommandWithOutput wrapper around commands returning their output and error
// NOTE: If you don't pass any formatArgs we'll just use the command directly,
// however there's a bizarre compiler error/warning when you pass in a formatString
// with a percent sign because it thinks it's supposed to be a formatString when
// in that case it's not. To get around that error you'll need to define the string
// in a variable and pass the variable into RunCommandWithOutput.
func (c *OSCommand) RunCommandWithOutput(formatString string, formatArgs ...interface{}) (string, error) {
	command := formatString
	if formatArgs != nil {
		command = fmt.Sprintf(formatString, formatArgs...)
	}
	cmd := c.ExecutableFromString(command)
	c.LogExecCmd(cmd)
	output, err := sanitisedCommandOutput(cmd.CombinedOutput())
	if err != nil {
		c.Log.WithField("command", command).Error(output)
	}
	return output, err
}

// RunExecutableWithOutput runs an executable file and returns its output
func (c *OSCommand) RunExecutableWithOutput(cmd *exec.Cmd) (string, error) {
	c.LogExecCmd(cmd)
	c.BeforeExecuteCmd(cmd)
	return sanitisedCommandOutput(cmd.CombinedOutput())
}

// RunExecutable runs an executable file and returns an error if there was one
func (c *OSCommand) RunExecutable(cmd *exec.Cmd) error {
	_, err := c.RunExecutableWithOutput(cmd)
	return err
}

// ExecutableFromString takes a string like `git status` and returns an executable command for it
func (c *OSCommand) ExecutableFromString(commandStr string) *exec.Cmd {
	splitCmd := str.ToArgv(commandStr)
	cmd := c.Command(splitCmd[0], splitCmd[1:]...)
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	return cmd
}

// ShellCommandFromString takes a string like `git commit` and returns an executable shell command for it
func (c *OSCommand) ShellCommandFromString(commandStr string) *exec.Cmd {
	quotedCommand := ""
	// Windows does not seem to like quotes around the command
	if c.Platform.OS == "windows" {
		quotedCommand = commandStr
	} else {
		quotedCommand = c.Quote(commandStr)
	}

	shellCommand := fmt.Sprintf("%s %s %s", c.Platform.Shell, c.Platform.ShellArg, quotedCommand)
	return c.ExecutableFromString(shellCommand)
}

// RunCommandWithOutputLive runs RunCommandWithOutputLiveWrapper
func (c *OSCommand) RunCommandWithOutputLive(command string, output func(string) string) error {
	return RunCommandWithOutputLiveWrapper(c, command, output)
}

func (c *OSCommand) CatFile(filename string) (string, error) {
	arr := append(c.Platform.CatCmd, filename)
	cmdStr := strings.Join(arr, " ")
	c.Log.WithField("command", cmdStr).Info("Cat")
	cmd := c.Command(arr[0], arr[1:]...)
	output, err := sanitisedCommandOutput(cmd.CombinedOutput())
	if err != nil {
		c.Log.WithField("command", cmdStr).Error(output)
	}
	return output, err
}

// DetectUnamePass detect a username / password / passphrase question in a command
// promptUserForCredential is a function that gets executed when this function detect you need to fillin a password or passphrase
// The promptUserForCredential argument will be "username", "password" or "passphrase" and expects the user's password/passphrase or username back
func (c *OSCommand) DetectUnamePass(command string, promptUserForCredential func(string) string) error {
	ttyText := ""
	errMessage := c.RunCommandWithOutputLive(command, func(word string) string {
		ttyText = ttyText + " " + word

		prompts := map[string]string{
			`.+'s password:`:                         "password",
			`Password\s*for\s*'.+':`:                 "password",
			`Username\s*for\s*'.+':`:                 "username",
			`Enter\s*passphrase\s*for\s*key\s*'.+':`: "passphrase",
		}

		for pattern, askFor := range prompts {
			if match, _ := regexp.MatchString(pattern, ttyText); match {
				ttyText = ""
				return promptUserForCredential(askFor)
			}
		}

		return ""
	})
	return errMessage
}

// RunCommand runs a command and just returns the error
func (c *OSCommand) RunCommand(formatString string, formatArgs ...interface{}) error {
	_, err := c.RunCommandWithOutput(formatString, formatArgs...)
	return err
}

// RunShellCommand runs shell commands i.e. 'sh -c <command>'. Good for when you
// need access to the shell
func (c *OSCommand) RunShellCommand(command string) error {
	cmd := c.Command(c.Platform.Shell, c.Platform.ShellArg, command)
	c.LogExecCmd(cmd)

	_, err := sanitisedCommandOutput(cmd.CombinedOutput())

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

// OpenFile opens a file with the given
func (c *OSCommand) OpenFile(filename string) error {
	commandTemplate := c.Config.GetUserConfig().OS.OpenCommand
	templateValues := map[string]string{
		"filename": c.Quote(filename),
	}

	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	err := c.RunCommand(command)
	return err
}

// OpenLink opens a file with the given
func (c *OSCommand) OpenLink(link string) error {
	c.LogCommand(fmt.Sprintf("Opening link '%s'", link), false)
	commandTemplate := c.Config.GetUserConfig().OS.OpenLinkCommand
	templateValues := map[string]string{
		"link": c.Quote(link),
	}

	command := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	err := c.RunCommand(command)
	return err
}

// PrepareSubProcess iniPrepareSubProcessrocess then tells the Gui to switch to it
// TODO: see if this needs to exist, given that ExecutableFromString does the same things
func (c *OSCommand) PrepareSubProcess(cmdName string, commandArgs ...string) *exec.Cmd {
	cmd := c.Command(cmdName, commandArgs...)
	if cmd != nil {
		cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	}
	c.LogExecCmd(cmd)
	return cmd
}

// PrepareShellSubProcess returns the pointer to a custom command
func (c *OSCommand) PrepareShellSubProcess(command string) *exec.Cmd {
	return c.PrepareSubProcess(c.Platform.Shell, c.Platform.ShellArg, command)
}

// Quote wraps a message in platform-specific quotation marks
func (c *OSCommand) Quote(message string) string {
	if c.Platform.OS == "windows" {
		message = strings.Replace(message, `"`, `"'"'"`, -1)
		message = strings.Replace(message, `\"`, `\\"`, -1)
	} else {
		message = strings.Replace(message, `\`, `\\`, -1)
		message = strings.Replace(message, `"`, `\"`, -1)
		message = strings.Replace(message, "`", "\\`", -1)
		message = strings.Replace(message, "$", "\\$", -1)
	}
	escapedQuote := c.Platform.EscapedQuote
	return escapedQuote + message + escapedQuote
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

// RunPreparedCommand takes a pointer to an exec.Cmd and runs it
// this is useful if you need to give your command some environment variables
// before running it
func (c *OSCommand) RunPreparedCommand(cmd *exec.Cmd) error {
	c.BeforeExecuteCmd(cmd)
	c.LogExecCmd(cmd)
	out, err := cmd.CombinedOutput()
	outString := string(out)
	c.Log.Info(outString)
	if err != nil {
		if len(outString) == 0 {
			return err
		}
		return errors.New(outString)
	}
	return nil
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
		cmds[i] = c.ExecutableFromString(str)
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

func RunLineOutputCmd(cmd *exec.Cmd, onLine func(line string) (bool, error)) error {
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

func (c *OSCommand) CopyToClipboard(str string) error {
	c.LogCommand(fmt.Sprintf("Copying '%s' to clipboard", utils.TruncateWithEllipsis(str, 40)), false)
	return clipboard.WriteAll(str)
}

func (c *OSCommand) RemoveFile(path string) error {
	c.LogCommand(fmt.Sprintf("Deleting path '%s'", path), false)

	return c.removeFile(path)
}
