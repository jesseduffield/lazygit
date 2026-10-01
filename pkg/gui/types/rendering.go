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

	// contentIsDiff marks output that is a panel's own diff; see ContentIsDiff.
	contentIsDiff bool
}

func (t *RunCommandTask) IsUpdateTask() {}

func NewRunCommandTask(cmd *exec.Cmd) *RunCommandTask {
	return &RunCommandTask{Cmd: cmd}
}

func NewRunCommandTaskWithPrefix(cmd *exec.Cmd, prefix string) *RunCommandTask {
	return &RunCommandTask{Cmd: cmd, Prefix: prefix}
}

type RunDiffRendererTask struct {
	Cmd    *exec.Cmd
	Prefix string

	// contentIsDiff marks output that is a panel's own diff; see ContentIsDiff.
	contentIsDiff bool
}

func (t *RunDiffRendererTask) IsUpdateTask() {}

func NewRunDiffRendererTask(cmd *exec.Cmd) *RunDiffRendererTask {
	return &RunDiffRendererTask{Cmd: cmd}
}

func NewRunDiffRendererTaskWithPrefix(cmd *exec.Cmd, prefix string) *RunDiffRendererTask {
	return &RunDiffRendererTask{Cmd: cmd, Prefix: prefix}
}

// NewMainViewDiffTask returns the task for rendering a diff into a main view. A diff
// normally goes through the diff renderer, however the render arranges to feed it. A
// diff we are producing with git itself, because the renderer's version of it couldn't
// be acted on, has to keep the renderer out, so it runs as a plain command instead.
//
// The task it returns is the one that says its output is a diff, so a pane rendering
// it can be pointed at (see ContentIsDiff).
func NewMainViewDiffTask(cmd *exec.Cmd, mode git_commands.DiffMode) UpdateTask {
	return NewMainViewDiffTaskWithPrefix(cmd, "", mode)
}

func NewMainViewDiffTaskWithPrefix(cmd *exec.Cmd, prefix string, mode git_commands.DiffMode) UpdateTask {
	if mode == git_commands.DiffModeRaw {
		task := NewRunCommandTaskWithPrefix(cmd, prefix)
		task.contentIsDiff = true
		return task
	}
	task := NewRunDiffRendererTaskWithPrefix(cmd, prefix)
	task.contentIsDiff = true
	return task
}

// ContentIsDiff reports whether the given render fills a main pane with the diff a
// panel offers there, as opposed to a message, a commit log, or a diff that is part of
// an explanation. A selection means the lines of the panel's diff, so it is only over
// such a render that there is anything to point at.
func ContentIsDiff(task UpdateTask) bool {
	switch task := task.(type) {
	case *RunCommandTask:
		return task.contentIsDiff
	case *RunDiffRendererTask:
		return task.contentIsDiff
	}
	return false
}
