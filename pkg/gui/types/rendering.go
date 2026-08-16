package types

import (
	"os/exec"
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
	Staging        MainContextPair
	PatchBuilding  MainContextPair
}

type ViewUpdateOpts struct {
	Title    string
	SubTitle string

	Task UpdateTask
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

type RunDiffRendererTask struct {
	Cmd    *exec.Cmd
	Prefix string
}

func (t *RunDiffRendererTask) IsUpdateTask() {}

func NewRunDiffRendererTask(cmd *exec.Cmd) *RunDiffRendererTask {
	return &RunDiffRendererTask{Cmd: cmd}
}

func NewRunDiffRendererTaskWithPrefix(cmd *exec.Cmd, prefix string) *RunDiffRendererTask {
	return &RunDiffRendererTask{Cmd: cmd, Prefix: prefix}
}
