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

func (self *DiffCommands) internalDiffCmdObj(diffArgs ...string) *GitCommandBuilder {
	return NewGitCmd("diff").
		Arg("--no-ext-diff", "--no-color").
		Arg(diffArgs...)
}

func (self *DiffCommands) GetPathDiff(path string, staged bool) (string, error) {
	return self.cmd.New(
		self.internalDiffCmdObj().
			ArgIf(staged, "--staged").
			Arg(path).
			ToArgv(),
	).RunWithOutput()
}

func (self *DiffCommands) GetAllDiff(staged bool) (string, error) {
	return self.cmd.New(
		self.internalDiffCmdObj().
			ArgIf(staged, "--staged").
			ToArgv(),
	).RunWithOutput()
}
