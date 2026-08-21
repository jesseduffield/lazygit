package git_commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/mgutz/str"
)

// metadataHandshake is what a diff renderer that speaks the OSC 1717 protocol emits
// before anything else, to announce that it does: a version-only record, with none of
// the fields a line's record has. See ProbeDiffRendererEmitsMetadata, and gocui's
// escape interpreter for how it is kept off the screen on a real render.
const metadataHandshake = "\x1b]1717"

// ProbeDiffRendererEmitsMetadata reports whether the configured diff renderer states
// which line of which file it is rendering, by running it on empty input and looking
// for the handshake. That is what decides whether a diff the renderer produced can be
// acted on at all, or has to be replaced by git's own when the user wants to act on it
// (see DiffLineHelper.MainViewDiffMode).
//
// Asking rather than watching a real render: the handshake is the renderer's first
// output whatever the diff, so the answer is a property of the renderer, known before
// we render anything — where watching would have to see a diff go by first, and would
// be fooled by a diff with no lines to describe.
//
// No terminal is needed. git only invokes a stdin filter when it thinks it is talking
// to one, but the renderer itself doesn't care: it announces itself whenever OSC1717 is
// set, so it can be run directly with empty input.
func (self *DiffCommands) ProbeDiffRendererEmitsMetadata() bool {
	manager := self.diffRendererConfigManager

	switch manager.GetDiffRendererType() {
	case config.DiffRendererType_StdinFilter:
		if command := manager.GetStdinFilterCommand(0); command != "" {
			return self.probeEmitsMetadata(self.cmd.NewShell(command, ""))
		}
	case config.DiffRendererType_ExtDiff:
		// An empty command means git's own diff.external config, which picks a driver
		// per file through .gitattributes: there is no one renderer to ask, and a single
		// diff can be produced by several, so we take it that it says nothing.
		if command := manager.GetExternalDiffCommand(3); command != "" {
			return self.externalDiffEmitsMetadata(command)
		}
	case config.DiffRendererType_RawGit:
		// git describes only the formats whose output can't be read back as a diff, and
		// asked with the renderer's own arguments it answers for exactly the format
		// those select: a handshake for a word diff, silence for a unified one. With no
		// arguments there is nothing to fall back to anyway, since this already is git's
		// own diff.
		if args := manager.GetRawGitArgs(); len(args) > 0 {
			return self.rawGitEmitsMetadata(args)
		}
	}

	return false
}

// rawGitEmitsMetadata asks git itself, run with the diff renderer's own arguments.
func (self *DiffCommands) rawGitEmitsMetadata(rawGitArgs []string) bool {
	oldPath, newPath, cleanup, ok := self.probeFiles()
	if !ok {
		return false
	}
	defer cleanup()

	return self.probeEmitsMetadata(self.cmd.New(
		NewGitCmd("diff").
			Arg("--no-index").
			Arg(rawGitArgs...).
			Arg(oldPath, newPath).
			ToArgv(),
	))
}

// externalDiffEmitsMetadata asks an external diff command, invoking it the way git
// invokes one — with the seven positional arguments of git's diff.external convention —
// over two empty files, so that it announces itself without having a diff to render.
func (self *DiffCommands) externalDiffEmitsMetadata(externalDiffCommand string) bool {
	oldPath, newPath, cleanup, ok := self.probeFiles()
	if !ok {
		return false
	}
	defer cleanup()

	args := append(str.ToArgv(externalDiffCommand),
		"probe", oldPath, "0000000", "100644", newPath, "0000000", "100644")
	return self.probeEmitsMetadata(self.cmd.New(args))
}

// probeFiles makes the two empty files a probe stands a diff up from, and the cleanup
// that removes them. Empty, because what the probe wants is for the renderer to announce
// itself, not for it to have anything to say.
func (self *DiffCommands) probeFiles() (string, string, func(), bool) {
	tempDir := self.os.GetTempDir()

	oldFile, err := os.CreateTemp(tempDir, "lazygit-probe-old-*")
	if err != nil {
		return "", "", nil, false
	}
	oldFile.Close()

	newFile, err := os.CreateTemp(tempDir, "lazygit-probe-new-*")
	if err != nil {
		os.Remove(oldFile.Name())
		return "", "", nil, false
	}
	newFile.Close()

	return oldFile.Name(), newFile.Name(), func() {
		os.Remove(oldFile.Name())
		os.Remove(newFile.Name())
	}, true
}

func (self *DiffCommands) probeEmitsMetadata(cmdObj *oscommands.CmdObj) bool {
	cmdObj.AddEnvVars("OSC1717=V1")
	// A renderer may well object to being handed nothing to render; what it said before
	// objecting is what we are after, and that is captured either way.
	output, _ := cmdObj.RunWithOutput()
	return strings.Contains(output, metadataHandshake)
}

type DiffCommands struct {
	*GitCommon
}

func NewDiffCommands(gitCommon *GitCommon) *DiffCommands {
	return &DiffCommands{
		GitCommon: gitCommon,
	}
}

// This is for generating diffs to be shown in the UI (e.g. rendering a range
// diff to the main view). It uses a custom diff renderer if one is configured.
func (self *DiffCommands) DiffCmdObj(diffArgs []string, mode DiffMode) *oscommands.CmdObj {
	return self.cmd.New(
		NewGitCmd("diff").
			Config("diff.noprefix=false").
			AddCommonDiffArgs(self.diffRendererConfigManager, self.UserConfig(), mode).
			Arg("--submodule").
			Arg(fmt.Sprintf("--color=%s", mode.colorArg(self.diffRendererConfigManager))).
			Arg(diffArgs...).
			Dir(self.repoPaths.worktreePath).
			ToArgv(),
	)
}

// CustomPatchDiffCmdObj is the command that renders the custom patch being built: a diff
// of the two trees the patch was materialized into (PatchCommands.WriteCustomPatchDiffTrees),
// under the directory holding them. It goes through the same wiring as any other diff we
// show, so the patch is rendered by whatever renders the rest of them, and git works out
// how much context to give it.
//
// git's own path prefixes are suppressed because the trees are named a and b themselves,
// which leaves the paths reading like an ordinary diff's over the repo's own paths.
func (self *DiffCommands) CustomPatchDiffCmdObj(dir string, mode DiffMode) *oscommands.CmdObj {
	return self.cmd.New(
		NewGitCmd("diff").
			AddCommonDiffArgs(self.diffRendererConfigManager, self.UserConfig(), mode).
			Arg("--no-index").
			Arg("--no-prefix").
			Arg(fmt.Sprintf("--color=%s", mode.colorArg(self.diffRendererConfigManager))).
			Arg("a", "b").
			Dir(dir).
			ToArgv(),
	)
}

// This is a basic generic diff command that can be used for any diff operation
// (e.g. copying a diff to the clipboard). It will not use a custom diff renderer,
// and does not use user configs such as ignore whitespace.
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
