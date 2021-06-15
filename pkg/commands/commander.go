package commands

import (
	"os"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

// to generate this run:
// counterfeiter pkg/commands ICommander
type ICommander interface {
	Run(cmdObj ICmdObj) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	RunGitCmdFromStr(cmdStr string) error
	BuildGitCmdObjFromStr(cmdStr string) ICmdObj
	BuildShellCmdObj(command string) ICmdObj
	SkipEditor(cmdObj ICmdObj)
	Quote(string) string
}

type runWithOutputFunc func(ICmdObj) (string, error)

type Commander struct {
	runWithOutput runWithOutputFunc
	log           *logrus.Entry
	lazygitPath   string
	shell         string // e.g. 'bash'
	shellArg      string // e.g. '-c'
	quote         func(string) string
}

func NewCommander(
	runWithOutput runWithOutputFunc,
	log *logrus.Entry,
	lazygitPath string,
	quote func(string) string,
) *Commander {
	return &Commander{
		runWithOutput: runWithOutput,
		log:           log,
		lazygitPath:   lazygitPath,
		quote:         quote,
	}
}

func (c *Commander) Run(cmdObj ICmdObj) error {
	_, err := c.RunWithOutput(cmdObj)
	return err
}

func (c *Commander) RunWithOutput(cmdObj ICmdObj) (string, error) {
	// TODO: have this retry logic in other places we run the command
	waitTime := 50 * time.Millisecond
	retryCount := 5
	attempt := 0

	for {
		output, err := c.runWithOutput(cmdObj)
		if err != nil {
			// if we have an error based on the index lock, we should wait a bit and then retry
			if strings.Contains(output, ".git/index.lock") {
				c.log.Error(output)
				c.log.Info("index.lock prevented command from running. Retrying command after a small wait")
				attempt++
				time.Sleep(waitTime)
				if attempt < retryCount {
					continue
				}
			}
		}
		return output, err
	}
}

func (c *Commander) RunGitCmdFromStr(cmdStr string) error {
	return c.Run(BuildGitCmdObjFromStr(cmdStr))
}

func BuildGitCmdStr(command string, positionalArgs []string, kwArgs map[string]bool) string {
	parts := []string{command}

	if len(kwArgs) > 0 {
		args := make([]string, 0, len(kwArgs))
		for arg, include := range kwArgs {
			if include {
				args = append(args, arg)
			}
		}
		utils.SortAlphabeticalInPlace(args)

		parts = append(parts, args...)
	}

	if len(positionalArgs) > 0 {
		parts = append(parts, positionalArgs...)
	}

	parts = utils.ExcludeEmpty(parts)

	return strings.Join(parts, " ")
}

func BuildGitCmdObj(command string, positionalArgs []string, kwArgs map[string]bool) ICmdObj {
	return BuildGitCmdObjFromStr(BuildGitCmdStr(command, positionalArgs, kwArgs))
}

// returns a command object from a command string. Prepends the `git ` part itself so
// if you want to do `git diff` just pass `diff` as the cmdStr
func BuildGitCmdObjFromStr(cmdStr string) ICmdObj {
	cmdObj := oscommands.NewCmdObjFromStr(GitCmdStr() + " " + cmdStr)
	SetDefaultEnvVars(cmdObj)

	return cmdObj
}

func BuildGitCmdObjFromArgs(args []string) ICmdObj {
	cmdObj := oscommands.NewCmdObjFromArgs(append([]string{GitCmdStr()}, args...))
	SetDefaultEnvVars(cmdObj)

	return cmdObj
}

func (c *Commander) BuildGitCmdObjFromStr(cmdStr string) ICmdObj {
	return BuildGitCmdObjFromStr(cmdStr)
}

func GitInitCmd() ICmdObj {
	return BuildGitCmdObjFromStr("init")
}

func GitVersionCmd() ICmdObj {
	return BuildGitCmdObjFromStr("--version")
}

func SetDefaultEnvVars(cmdObj ICmdObj) {
	cmdObj.GetCmd().Env = os.Environ()
	DisableOptionalLocks(cmdObj)
}

func DisableOptionalLocks(cmdObj ICmdObj) {
	cmdObj.AddEnvVars("GIT_OPTIONAL_LOCKS=0")
}

func (c *Commander) SkipEditor(cmdObj ICmdObj) {
	cmdObj.AddEnvVars(
		"LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY",
		"GIT_EDITOR="+c.lazygitPath,
		"EDITOR="+c.lazygitPath,
		"VISUAL="+c.lazygitPath,
	)
}

// We may have use for centralising this e.g. so that we call a specific executable
// or so that we can prepend some flags for bare repos.
// TODO: make this a method on the Git struct
func GitCmdStr() string {
	return "git"
}

// BuildShellCmdObj returns the pointer to a custom command
func (c *Commander) BuildShellCmdObj(command string) ICmdObj {
	return oscommands.NewCmdObjFromArgs([]string{c.shell, c.shellArg, command})
}

func (c *Commander) Quote(str string) string {
	return c.quote(str)
}
