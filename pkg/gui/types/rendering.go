package types

import (
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
)

type MainContextPair struct {
	Main      Context
	Secondary Context
}

func NewMainContextPair(main Context, secondary Context) MainContextPair {
	return MainContextPair{Main: main, Secondary: secondary}
}

// MainPanes says which of the two panes of the main section are shown. Most content
// takes the main pane alone; content with two sides to it — the working tree's
// unstaged and staged changes, a commit's diff and the patch built from it — takes
// both; and content whose only side is the second one takes the secondary pane alone,
// so that it has the whole section rather than sitting under an empty pane.
type MainPanes int

const (
	MainPaneOnly MainPanes = iota
	BothMainPanes
	SecondaryPaneOnly
)

type MainViewPairs struct {
	Normal         MainContextPair
	MergeConflicts MainContextPair
}

type ViewUpdateOpts struct {
	Title    string
	SubTitle string

	Task UpdateTask

	// NothingToActOn marks a pane that is being shown only because the layout is
	// configured to always split the diff: its side of the file holds nothing, so it
	// is not a pane to leave the focus in.
	NothingToActOn bool
}

type RefreshMainOpts struct {
	Pair      MainContextPair
	Main      *ViewUpdateOpts
	Secondary *ViewUpdateOpts
}

type UpdateTask interface {
	IsUpdateTask()
}

type RenderStringTask struct {
	Str string
}

func (t *RenderStringTask) IsUpdateTask() {}

func NewRenderStringTask(str string) *RenderStringTask {
	return &RenderStringTask{Str: str}
}

type RenderStringWithoutScrollTask struct {
	Str string
}

func (t *RenderStringWithoutScrollTask) IsUpdateTask() {}

func NewRenderStringWithoutScrollTask(str string) *RenderStringWithoutScrollTask {
	return &RenderStringWithoutScrollTask{Str: str}
}

type RenderStringWithScrollTask struct {
	Str     string
	OriginX int
	OriginY int
}

func (t *RenderStringWithScrollTask) IsUpdateTask() {}

func NewRenderStringWithScrollTask(str string, originX int, originY int) *RenderStringWithScrollTask {
	return &RenderStringWithScrollTask{Str: str, OriginX: originX, OriginY: originY}
}

type RunCommandTask struct {
	Cmd    *exec.Cmd
	Prefix string
}

func (t *RunCommandTask) IsUpdateTask() {}

func NewRunCommandTask(cmd *exec.Cmd) *RunCommandTask {
	return &RunCommandTask{Cmd: cmd}
}

func NewRunCommandTaskWithPrefix(cmd *exec.Cmd, prefix string) *RunCommandTask {
	return &RunCommandTask{Cmd: cmd, Prefix: prefix}
}

type RunPtyTask struct {
	Cmd    *exec.Cmd
	Prefix string
}

func (t *RunPtyTask) IsUpdateTask() {}

func NewRunPtyTask(cmd *exec.Cmd) *RunPtyTask {
	return &RunPtyTask{Cmd: cmd}
}

func NewRunPtyTaskWithPrefix(cmd *exec.Cmd, prefix string) *RunPtyTask {
	return &RunPtyTask{Cmd: cmd, Prefix: prefix}
}

// NewMainViewDiffTask returns the task for rendering a diff into a main view. Diffs
// normally run under a pty, since git only hands its output to a diff renderer when it
// thinks it is talking to a terminal — but a diff we are producing with git itself,
// because the renderer's version of it couldn't be acted on, has to keep the renderer
// out, so it runs as a plain command instead.
func NewMainViewDiffTask(cmd *exec.Cmd, mode git_commands.DiffMode) UpdateTask {
	return NewMainViewDiffTaskWithPrefix(cmd, "", mode)
}

func NewMainViewDiffTaskWithPrefix(cmd *exec.Cmd, prefix string, mode git_commands.DiffMode) UpdateTask {
	if mode == git_commands.DiffModeRaw {
		return NewRunCommandTaskWithPrefix(cmd, prefix)
	}
	return NewRunPtyTaskWithPrefix(cmd, prefix)
}
