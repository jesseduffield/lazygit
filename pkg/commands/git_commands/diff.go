package git_commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mgutz/str"
)

// metadataHandshake is the record a diff renderer that speaks the OSC 1717 protocol
// emits before anything else, to announce that it does: a version-only record, with
// none of the fields a line's record has. See ProbeDiffRendererEmitsMetadata, and
// gocui's escape interpreter for how it is kept off the screen on a real render.
const metadataHandshake = "\x1b]1717"

// probeWidth is the width a probe asks the renderer to lay out to. A probe has
// nothing to render, so any width does, as long as it is one a renderer will
// accept: a renderer that lays its rendering out in columns may well refuse a
// width of zero, and would then never get as far as announcing itself.
const probeWidth = 80

// ProbeDiffRendererEmitsMetadata reports whether the configured diff renderer states
// which line of which file it is rendering, by running it on empty input and looking
// for the handshake. The answer decides whether a diff the renderer produced can be
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
	values := config.DiffRendererValues{Width: probeWidth, DiffContext: 3}

	switch manager.GetDiffRendererType() {
	case config.DiffRendererType_StdinFilter:
		if command, err := manager.GetStdinFilterCommand(values); err == nil && command != "" {
			return self.probeEmitsMetadata(self.cmd.NewShell(command, ""))
		}
	case config.DiffRendererType_ExtDiff:
		// An empty command means git's own diff.external config, which picks a driver
		// per file through .gitattributes: there is no one renderer to ask, and a single
		// diff can be produced by several, so we take it that it says nothing.
		if command, err := manager.GetExternalDiffCommand(values); err == nil && command != "" {
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
	// A renderer may well object to being handed nothing to render. We want to know
	// whatever it said before objecting, and that is captured either way.
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
			NoLineEndingConversion().
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

// NamedText is a text to be diffed by a TextDiffer, along with the name to show
// for it in the diff. The name doubles as the name of the temporary file holding
// the text, so it has to be usable as a file name.
type NamedText struct {
	Name    string
	Content string
}

// A TextDiffer diffs two texts with the configured diff renderer, in the same way
// as every other diff we show. It holds everything it takes from the
// configuration, so that it can render a diff on any goroutine.
type TextDiffer struct {
	diffCommands *DiffCommands

	// The git command that diffs the two texts, up to the names of the files
	// holding them
	args         []string
	rawGit       bool
	stdinFilter  string
	externalDiff string
	width        int
	height       int
}

// NewTextDiffer resolves what a TextDiffer takes from the configuration. values
// are what the renderer's command is resolved with. values.Width and height are
// the size of the view the diff is going to be shown in; a diff renderer lays its
// output out for them.
//
// It has to be called on the UI thread, since that is where the user switches
// diff renderers and has the config reloaded.
func (self *DiffCommands) NewTextDiffer(values config.DiffRendererValues, height int) (*TextDiffer, error) {
	manager := self.diffRendererConfigManager
	stdinFilter, err := manager.GetStdinFilterCommand(values)
	if err != nil {
		return nil, err
	}
	externalDiff, err := manager.GetExternalDiffCommand(values)
	if err != nil {
		return nil, err
	}

	return &TextDiffer{
		diffCommands: self,
		args: NewGitCmd("diff").
			// git heads a hunk with the nearest line above it that starts with
			// a letter, for the name of the function the hunk is in. In a text,
			// that is just a line out of context. A pattern that never matches
			// leaves it out; "default" is the driver of a file that no
			// attribute gives one.
			Config("diff.default.xfuncname=x^").
			AddCommonDiffArgs(manager, self.UserConfig(), DiffModeRendered).
			Arg("--no-index", "--no-prefix").
			Arg(fmt.Sprintf("--color=%s", DiffModeRendered.colorArg(manager))).
			ToArgv(),
		rawGit:       manager.GetDiffRendererType() == config.DiffRendererType_RawGit,
		stdinFilter:  stdinFilter,
		externalDiff: externalDiff,
		width:        values.Width,
		height:       height,
	}, nil
}

// RenderedDiff returns a diff of the two given texts.
//
// git can only diff files, so the texts are written to temporary files named
// after them. Under a diff renderer those names are what the diff calls the two
// sides; under git's own diff they go away with the rest of the header.
func (self *TextDiffer) RenderedDiff(before NamedText, after NamedText) (string, error) {
	dir, err := os.MkdirTemp(self.diffCommands.os.GetTempDir(), "textdiff-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	for _, text := range []NamedText{before, after} {
		content := text.Content
		// End the file with a newline, or git says that it doesn't, which tells
		// a reader of the diff nothing about the two texts.
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		if err := self.diffCommands.os.CreateFileWithContent(filepath.Join(dir, text.Name), content); err != nil {
			return "", err
		}
	}

	cmdObj := self.diffCommands.cmd.New(
		append(slices.Clone(self.args), "--", before.Name, after.Name),
	).SetWd(dir).DontLog()

	// --no-index implies --exit-code, so git exits with a non-zero status
	// whenever the two texts differ; that is only an error if it left us with
	// nothing to show.
	if self.rawGit {
		// git renders the diff itself here, so no terminal is needed to get it.
		output, err := cmdObj.RunWithOutput()
		if output == "" && err != nil {
			return "", err
		}

		return stripDiffHeaders(output), nil
	}

	oscommands.SetDumbTerminalEnv(cmdObj.GetCmd())
	// For diff renderers that don't ask the terminal how wide it is.
	cmdObj.AddEnvVars(fmt.Sprintf("COLUMNS=%d", self.width))
	// An empty command means that git's own diff.external config applies, and
	// git consults that only while the variable is unset.
	if self.externalDiff != "" {
		cmdObj.AddEnvVars("GIT_EXTERNAL_DIFF=" + self.externalDiff)
	}

	var output string
	if oscommands.RendersThroughAPipe() {
		output, err = self.runThroughAPipe(cmdObj)
	} else {
		// git runs the stdin filter itself, as the pager it is told about here.
		// Named even when there is none, so that git doesn't reach for the
		// user's core.pager instead.
		cmdObj.AddEnvVars("GIT_PAGER=" + self.stdinFilter)
		output, err = oscommands.RunInPtyWithOutput(cmdObj.GetCmd(), uint16(self.width), uint16(self.height))
	}
	if output == "" && err != nil {
		return "", err
	}

	return output, nil
}

// runThroughAPipe feeds the output of cmdObj to the stdin filter through a pipe
// and returns what the filter wrote. git only runs a stdin filter itself when it
// thinks it is talking to a terminal, so here the filter is a command of our
// own. An external diff renderer is run by git, so with one the command runs
// alone.
func (self *TextDiffer) runThroughAPipe(cmdObj *oscommands.CmdObj) (string, error) {
	if self.stdinFilter == "" {
		return cmdObj.RunWithOutput()
	}

	// The filter runs in a plain shell, with git's environment, since that is
	// how git would have run it.
	pipeline, reader, err := self.diffCommands.os.StartPipeline(
		cmdObj,
		self.diffCommands.cmd.NewShell(self.stdinFilter, "").SetEnviron(cmdObj.GetCmd().Env).DontLog(),
	)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// The output ends once every command in the pipeline has exited, and Wait
	// says how they did.
	output, _ := io.ReadAll(reader)
	return string(output), pipeline.Wait()
}

// Strips the file header and the first hunk header from a diff. For a diff of
// two temporary files these say nothing that a reader could make sense of.
func stripDiffHeaders(diff string) string {
	lines := strings.SplitAfter(diff, "\n")
	for i, line := range lines {
		if strings.HasPrefix(utils.Decolorise(line), "@@ ") {
			return strings.Join(lines[i+1:], "")
		}
	}

	return diff
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
