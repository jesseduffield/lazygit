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

// This is for generating diffs to be shown in the UI (e.g. rendering a range
// diff to the main view). It uses a custom pager if one is configured.
func (self *DiffCommands) DiffCmdObj(diffArgs []string) oscommands.ICmdObj {
	extDiffCmd := self.UserConfig().Git.Paging.ExternalDiffCommand
	useExtDiff := extDiffCmd != ""
	ignoreWhitespace := self.AppState.IgnoreWhitespaceInDiffView

	return self.cmd.New(
		NewGitCmd("diff").
			Config("diff.noprefix=false").
			ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
			ArgIfElse(useExtDiff, "--ext-diff", "--no-ext-diff").
			Arg("--submodule").
			Arg(fmt.Sprintf("--color=%s", self.UserConfig().Git.Paging.ColorArg)).
			ArgIf(ignoreWhitespace, "--ignore-all-space").
			Arg(fmt.Sprintf("--unified=%d", self.AppState.DiffContextSize)).
			Arg(diffArgs...).
			Dir(self.repoPaths.worktreePath).
			ToArgv(),
	)
}

// This is a basic generic diff command that can be used for any diff operation
// (e.g. copying a diff to the clipboard). It will not use a custom pager, and
// does not use user configs such as ignore whitespace.
// If you want to diff specific refs (one or two), you need to add them yourself
// in additionalArgs; it is recommended to also pass `--` after that. If you
// want to restrict the diff to specific paths, pass them in additionalArgs
// after the `--`.
func (self *DiffCommands) GetDiff(staged bool, additionalArgs ...string) (string, error) {
	return self.cmd.New(
		NewGitCmd("diff").
			Config("diff.noprefix=false").
			Arg("--no-ext-diff", "--no-color").
			ArgIf(staged, "--staged").
			Dir(self.repoPaths.worktreePath).
			Arg(additionalArgs...).
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
