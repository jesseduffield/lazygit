package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/sirupsen/logrus"
)

// all we're doing here is wrapping the default command object builder with
// some git-specific stuff: e.g. adding a git-specific env var

type gitCmdObjBuilder struct {
	innerBuilder *oscommands.CmdObjBuilder
}

var _ oscommands.ICmdObjBuilder = &gitCmdObjBuilder{}

// We disable git's optional locks on every command by default so that our git
// invocations never contend for index.lock. See git_commands.OptionalLocksEnvVar
// for the full rationale. Individual commands that do want the lock (currently
// only the foreground files refresh) opt back in via CmdObj.RemoveEnvVar.
var defaultEnvVar = git_commands.OptionalLocksEnvVar + "=0"

func NewGitCmdObjBuilder(log *logrus.Entry, innerBuilder *oscommands.CmdObjBuilder) *gitCmdObjBuilder {
	// the price of having a convenient interface where we can say .New(...).Run() is that our builder now depends on our runner, so when we want to wrap the default builder/runner in new functionality we need to jump through some hoops. We could avoid the use of a decorator function here by just exporting the runner field on the default builder but that would be misleading because we don't want anybody using that to run commands (i.e. we want there to be a single API used across the codebase)
	updatedBuilder := innerBuilder.CloneWithNewRunner(func(runner oscommands.ICmdObjRunner) oscommands.ICmdObjRunner {
		return &gitCmdObjRunner{
			log:               log,
			innerRunner:       runner,
			initialRetryDelay: defaultInitialRetryDelay,
		}
	})

	return &gitCmdObjBuilder{
		innerBuilder: updatedBuilder,
	}
}

func (self *gitCmdObjBuilder) New(args []string) *oscommands.CmdObj {
	return self.innerBuilder.New(args).AddEnvVars(defaultEnvVar)
}

func (self *gitCmdObjBuilder) NewShell(cmdStr string, shellFunctionsFile string) *oscommands.CmdObj {
	return self.innerBuilder.NewShell(cmdStr, shellFunctionsFile).AddEnvVars(defaultEnvVar)
}

func (self *gitCmdObjBuilder) Quote(str string) string {
	return self.innerBuilder.Quote(str)
}
