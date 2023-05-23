package git_commands

import "github.com/jesseduffield/lazygit/pkg/commands/oscommands"

type DiffCommands struct {
	*GitCommon
}

func NewDiffCommands(gitCommon *GitCommon) *DiffCommands {
	return &DiffCommands{
		GitCommon: gitCommon,
	}
}

func (self *DiffCommands) DiffCmdObj(diffArgs []string) oscommands.ICmdObj {
	return self.cmd.New(
		NewGitCmd("diff").Arg("--submodule", "--no-ext-diff", "--color").Arg(diffArgs...).ToArgv(),
	)
}
