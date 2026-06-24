package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/sirupsen/logrus"
)

// NewGitCmdObjBuilder returns a command object builder whose runner is wrapped
// with our git-specific runner (logging, credential handling, etc.).
func NewGitCmdObjBuilder(log *logrus.Entry, innerBuilder *oscommands.CmdObjBuilder) *oscommands.CmdObjBuilder {
	// We decorate the runner rather than exposing the builder's runner field:
	// that field stays unexported so there's a single API for running commands
	// across the codebase.
	return innerBuilder.CloneWithNewRunner(func(runner oscommands.ICmdObjRunner) oscommands.ICmdObjRunner {
		return &gitCmdObjRunner{
			log:         log,
			innerRunner: runner,
		}
	})
}
