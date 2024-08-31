package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/sirupsen/logrus"
)

// all we're doing here is wrapping the default command object builder with
// some git-specific stuff: e.g. adding a git-specific env var

type gitCmdObjBuilder struct {
	innerBuilder *oscommands.CmdObjBuilder
}

var _ oscommands.ICmdObjBuilder = &gitCmdObjBuilder{}

func NewGitCmdObjBuilder(log *logrus.Entry, innerBuilder *oscommands.CmdObjBuilder) *gitCmdObjBuilder {
	// the price of having a convenient interface where we can say .New(...).Run() is that our builder now depends on our runner, so when we want to wrap the default builder/runner in new functionality we need to jump through some hoops. We could avoid the use of a decorator function here by just exporting the runner field on the default builder but that would be misleading because we don't want anybody using that to run commands (i.e. we want there to be a single API used across the codebase)
	updatedBuilder := innerBuilder.CloneWithNewRunner(func(runner oscommands.ICmdObjRunner) oscommands.ICmdObjRunner {
		return &gitCmdObjRunner{
			log:         log,
			innerRunner: runner,
		}
	})

	return &gitCmdObjBuilder{
		innerBuilder: updatedBuilder,
	}
}

var defaultEnvVar = "GIT_OPTIONAL_LOCKS=0"

func (self *gitCmdObjBuilder) New(args []string) oscommands.ICmdObj {
	return self.innerBuilder.New(args).AddEnvVars(defaultEnvVar)
}

func (self *gitCmdObjBuilder) NewShell(cmdStr string) oscommands.ICmdObj {
	return self.innerBuilder.NewShell(cmdStr).AddEnvVars(defaultEnvVar)
}

func (self *gitCmdObjBuilder) NewInteractiveShell(cmdStr string) oscommands.ICmdObj {
	return self.innerBuilder.NewInteractiveShell(cmdStr).AddEnvVars(defaultEnvVar)
}

func (self *gitCmdObjBuilder) Quote(str string) string {
	return self.innerBuilder.Quote(str)
}
