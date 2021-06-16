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
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

type IOS interface {
	Getenv(envVar string) string
	WithSpan(span string) IOS
	LogCmd(cmd ICmdObj)
	LogCommand(cmdStr string, commandLine bool)
	SetOnRunCommand(f func(CmdLogEntry))
	SetCommand(cmd func(string, ...string) *exec.Cmd)
	SetRemoveFile(f func(string) error)
	CatFile(filename string) (string, error)
	FileType(path string) string
	OpenFile(filename string) error
	OpenLink(link string) error
	Quote(message string) string
	AppendLineToFile(filename, line string) error
	CreateTempFile(filename, content string) (string, error)
	CreateFileWithContent(path string, content string) error
	Remove(filename string) error
	FileExists(path string) (bool, error)
	GetLazygitPath() string
	PipeCommands(cmdObjs ...ICmdObj) error
	CopyToClipboard(str string) error
	RemoveFile(path string) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	Run(cmd ICmdObj) error
	RunAndParseWords(cmdObj ICmdObj, output func(string) string) error
}

// OS holds all the os commands
type OS struct {
	log      *logrus.Entry
	platform *Platform
	config   config.AppConfigurer
	command  func(string, ...string) *exec.Cmd
	getenv   func(string) string

	// callback to run before running a command, i.e. for the purposes of logging
	onRunCommand func(CmdLogEntry)

	// something like 'Staging File': allows us to group cmd logs under a single title
	cmdLogSpan string

	removeFile func(string) error
}

// NewOS os command runner
func NewOS(log *logrus.Entry, config config.AppConfigurer) *OS {
	return &OS{
		log:        log,
		platform:   getPlatform(),
		config:     config,
		command:    secureexec.Command,
		getenv:     os.Getenv,
		removeFile: os.RemoveAll,
	}
}

func (c *OS) Getenv(envVar string) string {
	return c.getenv(envVar)
}

func (c *OS) WithSpan(span string) IOS {
	// sometimes .WithSpan(span) will be called where span actually is empty, in
	// which case we don't need to log anything so we can just return early here
	// with the original struct
	if span == "" {
		return c
	}

	newOS := &OS{}
	*newOS = *c
	newOS.cmdLogSpan = span
	return newOS
}

func (c *OS) LogCmd(cmd ICmdObj) {
	c.LogCommand(cmd.ToString(), true)
}

func (c *OS) LogCommand(cmdStr string, commandLine bool) {
	c.log.WithField("command", cmdStr).Info("RunCommand")

	if c.onRunCommand != nil && c.cmdLogSpan != "" {
		c.onRunCommand(NewCmdLogEntry(cmdStr, c.cmdLogSpan, commandLine))
	}
}

func (c *OS) SetOnRunCommand(f func(CmdLogEntry)) {
	c.onRunCommand = f
}

// SetCommand sets the command function used by the struct.
// To be used for testing only
func (c *OS) SetCommand(cmd func(string, ...string) *exec.Cmd) {
	c.command = cmd
}

// To be used for testing only
func (c *OS) SetRemoveFile(f func(string) error) {
	c.removeFile = f
}

func (c *OS) CatFile(filename string) (string, error) {
	cmdObj := NewCmdObjFromArgs(append(c.platform.CatCmd, filename))
	c.log.WithField("command", cmdObj.ToString()).Info("Cat")
	return c.RunWithOutput(cmdObj)
}

// FileType tells us if the file is a file, directory or other
func (c *OS) FileType(path string) string {
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
func (c *OS) OpenFile(filename string) error {
	commandTemplate := c.config.GetUserConfig().OS.OpenCommand
	templateValues := map[string]string{
		"filename": c.Quote(filename),
	}

	cmdStr := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	return c.Run(NewCmdObjFromStr(cmdStr))
}

// OpenLink opens a file with the given
func (c *OS) OpenLink(link string) error {
	c.LogCommand(fmt.Sprintf("Opening link '%s'", link), false)
	commandTemplate := c.config.GetUserConfig().OS.OpenLinkCommand
	templateValues := map[string]string{
		"link": c.Quote(link),
	}

	cmdStr := utils.ResolvePlaceholderString(commandTemplate, templateValues)
	return c.Run(NewCmdObjFromStr(cmdStr))
}

// Quote wraps a message in platform-specific quotation marks
func (c *OS) Quote(message string) string {
	if c.platform.OS == "windows" {
		message = strings.Replace(message, `"`, `"'"'"`, -1)
		message = strings.Replace(message, `\"`, `\\"`, -1)
	} else {
		message = strings.Replace(message, `\`, `\\`, -1)
		message = strings.Replace(message, `"`, `\"`, -1)
		message = strings.Replace(message, "`", "\\`", -1)
		message = strings.Replace(message, "$", "\\$", -1)
	}
	escapedQuote := c.platform.EscapedQuote
	return escapedQuote + message + escapedQuote
}

// AppendLineToFile adds a new line in file
func (c *OS) AppendLineToFile(filename, line string) error {
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
func (c *OS) CreateTempFile(filename, content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		c.log.Error(err)
		return "", utils.WrapError(err)
	}
	c.LogCommand(fmt.Sprintf("Creating temp file '%s'", tmpfile.Name()), false)

	if _, err := tmpfile.WriteString(content); err != nil {
		c.log.Error(err)
		return "", utils.WrapError(err)
	}
	if err := tmpfile.Close(); err != nil {
		c.log.Error(err)
		return "", utils.WrapError(err)
	}

	return tmpfile.Name(), nil
}

// CreateFileWithContent creates a file with the given content
func (c *OS) CreateFileWithContent(path string, content string) error {
	c.LogCommand(fmt.Sprintf("Creating file '%s'", path), false)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		c.log.Error(err)
		return err
	}

	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		c.log.Error(err)
		return utils.WrapError(err)
	}

	return nil
}

// Remove removes a file or directory at the specified path
func (c *OS) Remove(filename string) error {
	c.LogCommand(fmt.Sprintf("Removing '%s'", filename), false)
	err := os.RemoveAll(filename)
	return utils.WrapError(err)
}

// FileExists checks whether a file exists at the specified path
func (c *OS) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetLazygitPath returns the path of the currently executed file
func (c *OS) GetLazygitPath() string {
	ex, err := os.Executable() // get the executable path for git to use
	if err != nil {
		ex = os.Args[0] // fallback to the first call argument if needed
	}
	return `"` + filepath.ToSlash(ex) + `"`
}

// PipeCommands runs a heap of commands and pipes their inputs/outputs together like A | B | C
func (c *OS) PipeCommands(cmdObjs ...ICmdObj) error {
	logCmdStr := ""
	for i, cmdObj := range cmdObjs {
		if i > 0 {
			logCmdStr += " | "
		}
		logCmdStr += cmdObj.ToString()
	}
	c.LogCommand(logCmdStr, true)

	for i := 0; i < len(cmdObjs)-1; i++ {
		stdout, err := cmdObjs[i].GetCmd().StdoutPipe()
		if err != nil {
			return err
		}

		cmdObjs[i+1].GetCmd().Stdin = stdout
	}

	// keeping this here in case I adapt this code for some other purpose in the future
	// cmds[len(cmds)-1].Stdout = os.Stdout

	finalErrors := []string{}

	wg := sync.WaitGroup{}
	wg.Add(len(cmdObjs))

	for _, cmdObj := range cmdObjs {
		currentCmd := cmdObj.GetCmd()
		go utils.Safe(func() {
			stderr, err := currentCmd.StderrPipe()
			if err != nil {
				c.log.Error(err)
			}

			if err := currentCmd.Start(); err != nil {
				c.log.Error(err)
			}

			if b, err := ioutil.ReadAll(stderr); err == nil {
				if len(b) > 0 {
					finalErrors = append(finalErrors, string(b))
				}
			}

			if err := currentCmd.Wait(); err != nil {
				c.log.Error(err)
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

func (c *OS) CopyToClipboard(str string) error {
	c.LogCommand(fmt.Sprintf("Copying '%s' to clipboard", utils.TruncateWithEllipsis(str, 40)), false)
	return clipboard.WriteAll(str)
}

func (c *OS) RemoveFile(path string) error {
	c.LogCommand(fmt.Sprintf("Deleting path '%s'", path), false)

	return c.removeFile(path)
}
