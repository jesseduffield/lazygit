package git_commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/mgutz/str"
)

type DiffCommands struct {
	*GitCommon
}

func NewDiffCommands(gitCommon *GitCommon) *DiffCommands {
	return &DiffCommands{
		GitCommon: gitCommon,
	}
}

// metadataHandshake is the OSC sequence a metadata-aware pager emits, as its first
// output, to announce that it speaks the diff-line-metadata protocol: a version-only
// OSC 1717 record (no fields). See ProbePagerEmitsDiffMetadata and, for how it's
// swallowed on a real render, escapeInterpreter.dropMetadataIfHandshake.
const metadataHandshake = "\x1b]1717"

// ProbePagerEmitsDiffMetadata reports whether the configured pager speaks the
// diff-line-metadata protocol, by running it on empty input and checking for its
// handshake. It's the focused main view's signal for whether it can act on the
// pager's rendered diff or must fall back to the raw diff (see
// StagingHelper.DiffMainViewShouldRenderRaw). The verdict is content-independent —
// the handshake is the pager's first output regardless of the diff — so the caller
// caches it per pager.
//
// No PTY is needed: git needs a terminal to decide to invoke a pager, but the pager
// itself emits the handshake whenever OSC1717 is set, so we can run it
// directly with empty input.
//
// A git-config external diff driver (useExternalDiffGitConfig) is chosen per file via
// .gitattributes and a single diff can mix drivers, so there's no one pager to probe;
// we conservatively report false (the focused main view then always renders raw).
func (self *DiffCommands) ProbePagerEmitsDiffMetadata() bool {
	if self.pagerConfig.GetUseExternalDiffGitConfig() {
		return false
	}
	if extDiffCmd := self.pagerConfig.GetExternalDiffCommand(); extDiffCmd != "" {
		return self.externalDiffEmitsMetadata(extDiffCmd)
	}
	if pagerCmd := self.pagerConfig.GetPagerCommand(0); pagerCmd != "" {
		return self.probeEmitsMetadata(self.cmd.NewShell(pagerCmd, ""))
	}
	return false
}

// externalDiffEmitsMetadata probes an external diff command, invoking it the way git
// invokes a diff.external driver — with 7 positional args
// (path old-file old-hex old-mode new-file new-hex new-mode) — but on two empty temp
// files, so it emits its handshake without there being a real diff to render.
func (self *DiffCommands) externalDiffEmitsMetadata(extDiffCmd string) bool {
	oldFile, err := os.CreateTemp("", "lazygit-probe-old-*")
	if err != nil {
		return false
	}
	defer os.Remove(oldFile.Name())
	oldFile.Close()

	newFile, err := os.CreateTemp("", "lazygit-probe-new-*")
	if err != nil {
		return false
	}
	defer os.Remove(newFile.Name())
	newFile.Close()

	args := append(str.ToArgv(extDiffCmd),
		"probe", oldFile.Name(), "0000000", "100644", newFile.Name(), "0000000", "100644")
	return self.probeEmitsMetadata(self.cmd.New(args))
}

func (self *DiffCommands) probeEmitsMetadata(cmdObj *oscommands.CmdObj) bool {
	cmdObj.AddEnvVars("OSC1717=V1")
	// The pager may exit non-zero on the synthetic input; we only care about whether
	// it emitted the handshake first, and the output is captured either way.
	output, _ := cmdObj.RunWithOutput()
	return strings.Contains(output, metadataHandshake)
}

// This is for generating diffs to be shown in the UI (e.g. rendering a range
// diff to the main view). It uses a custom pager if one is configured, unless
// ignoreExternalDiff is set (the focused main view's raw-diff fallback; keeps the
// colour, unlike a plain diff).
func (self *DiffCommands) DiffCmdObj(diffArgs []string, ignoreExternalDiff bool) *oscommands.CmdObj {
	extDiffCmd := self.pagerConfig.GetExternalDiffCommand()
	useExtDiff := extDiffCmd != "" && !ignoreExternalDiff
	useExtDiffGitConfig := self.pagerConfig.GetUseExternalDiffGitConfig() && !ignoreExternalDiff
	ignoreWhitespace := self.UserConfig().Git.IgnoreWhitespaceInDiffView

	return self.cmd.New(
		NewGitCmd("diff").
			Config("diff.noprefix=false").
			ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
			ArgIfElse(useExtDiff || useExtDiffGitConfig, "--ext-diff", "--no-ext-diff").
			Arg("--submodule").
			Arg(fmt.Sprintf("--color=%s", self.pagerConfig.GetColorArg())).
			ArgIf(ignoreWhitespace, "--ignore-all-space").
			Arg(fmt.Sprintf("--unified=%d", self.UserConfig().Git.DiffContextSize)).
			Arg(diffArgs...).
			Dir(self.repoPaths.worktreePath).
			ToArgv(),
	)
}

// CustomPatchDiffCmdObj builds the command that renders the custom patch shown in the
// secondary pane: a `git diff --no-index` of the two file trees PatchCommands materialized
// under dir (a/ = before, b/ = after; see WriteCustomPatchDiffTrees). It uses the same pager
// wiring as DiffCmdObj so the patch renders exactly like any other diff — through a stdin
// pager, an external diff tool, or (when ignoreExternalDiff is set, the focused main view's
// raw-diff fallback) git's own colour. --no-prefix is used because the a/ and b/ tree names
// already stand in for git's conventional a//b/ path prefixes, so the diff's paths come out
// as the real repo-relative paths.
func (self *DiffCommands) CustomPatchDiffCmdObj(dir string, ignoreExternalDiff bool) *oscommands.CmdObj {
	extDiffCmd := self.pagerConfig.GetExternalDiffCommand()
	useExtDiff := extDiffCmd != "" && !ignoreExternalDiff
	useExtDiffGitConfig := self.pagerConfig.GetUseExternalDiffGitConfig() && !ignoreExternalDiff

	return self.cmd.New(
		NewGitCmd("diff").
			ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
			ArgIfElse(useExtDiff || useExtDiffGitConfig, "--ext-diff", "--no-ext-diff").
			Arg("--no-index").
			Arg("--no-prefix").
			Arg(fmt.Sprintf("--color=%s", self.pagerConfig.GetColorArg())).
			Arg(fmt.Sprintf("--unified=%d", self.UserConfig().Git.DiffContextSize)).
			Arg("a", "b").
			Dir(dir).
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
	).DontLog().RunWithOutput()
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

func (self *DiffCommands) OpenDiffToolCmdObj(opts DiffToolCmdOptions) *oscommands.CmdObj {
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

func (self *DiffCommands) DiffIndexCmdObj(diffArgs ...string) *oscommands.CmdObj {
	return self.cmd.New(
		NewGitCmd("diff-index").
			Config("diff.noprefix=false").
			Arg("--submodule", "--no-ext-diff", "--no-color", "--patch").
			Arg(diffArgs...).ToArgv(),
	)
}
