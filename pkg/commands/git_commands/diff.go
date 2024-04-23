package git_commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type DiffCommands struct {
	*GitCommon
}

func NewDiffCommands(gitCommon *GitCommon) *DiffCommands {
	return &DiffCommands{
		GitCommon: gitCommon,
	}
}

func (self *DiffCommands) DiffCmdObj(diffArgs []string) oscommands.ICmdObj {
	extDiffCmd := self.UserConfig.Git.Paging.ExternalDiffCommand
	useExtDiff := extDiffCmd != ""

	return self.cmd.New(
		NewGitCmd("diff").
			Config("diff.noprefix=false").
			ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
			ArgIfElse(useExtDiff, "--ext-diff", "--no-ext-diff").
			Arg("--submodule").
			Arg(fmt.Sprintf("--color=%s", self.UserConfig.Git.Paging.ColorArg)).
			Arg(diffArgs...).
			Dir(self.repoPaths.worktreePath).
			ToArgv(),
	)
}

func (self *DiffCommands) internalDiffCmdObj(diffArgs ...string) *GitCommandBuilder {
	return NewGitCmd("diff").
		Config("diff.noprefix=false").
		Arg("--no-ext-diff", "--no-color").
		Arg(diffArgs...).
		Dir(self.repoPaths.worktreePath)
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

type DiffToolCmdOptions struct {
	// The path to show a diff for. Pass "." for the entire repo.
	Filepath string

	// The commit against which to show the diff. Leave empty to show a diff of
	// the working copy.
	FromCommit string

	// The commit to diff against FromCommit. Leave empty to diff the working
	// copy against FromCommit. Leave both FromCommit and ToCommit empty to show
	// the diff of the unstaged working copy changes against the index if Staged
	// is false, or the staged changes against HEAD if Staged is true.
	ToCommit string

	// Whether to reverse the left and right sides of the diff.
	Reverse bool

	// Whether the given Filepath is a directory. We'll pass --dir-diff to
	// git-difftool in that case.
	IsDirectory bool

	// Whether to show the staged or the unstaged changes. Must be false if both
	// FromCommit and ToCommit are non-empty.
	Staged bool
}

func (self *DiffCommands) OpenDiffToolCmdObj(opts DiffToolCmdOptions) oscommands.ICmdObj {
	return self.cmd.New(NewGitCmd("difftool").
		Arg("--no-prompt").
		ArgIf(opts.IsDirectory, "--dir-diff").
		ArgIf(opts.Staged, "--cached").
		ArgIf(opts.FromCommit != "", opts.FromCommit).
		ArgIf(opts.ToCommit != "", opts.ToCommit).
		ArgIf(opts.Reverse, "-R").
		Arg("--", opts.Filepath).
		ToArgv())
}

func (self *DiffCommands) DiffIndexCmdObj(diffArgs ...string) oscommands.ICmdObj {
	return self.cmd.New(
		NewGitCmd("diff-index").
			Config("diff.noprefix=false").
			Arg("--submodule", "--no-ext-diff", "--no-color", "--patch").
			Arg(diffArgs...).ToArgv(),
	)
}
