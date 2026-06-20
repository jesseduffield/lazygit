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

// NewMainViewDiffTask builds the task for rendering a diff into the main view,
// choosing between the normal pty task and the focused main view's raw-diff
// fallback. When renderRaw is set (the focused main view needs to act on a diff the
// configured pager can't resolve) it uses a plain command task, which — unlike the
// pty task — doesn't pipe the diff through a stdin pager (GIT_PAGER); the external
// diff command, if any, is suppressed in the cmd itself (its ignoreExternalDiff arg).
// The caller passes the same renderRaw to the diff-cmd builder so the two stay in step.
func NewMainViewDiffTask(renderRaw bool, cmd *exec.Cmd) UpdateTask {
	if renderRaw {
		return NewRunCommandTask(cmd)
	}
	return NewRunPtyTask(cmd)
}

// NewMainViewDiffTaskWithPrefix is NewMainViewDiffTask for a diff rendered with a
// leading prefix (e.g. a range-diff or stash header).
func NewMainViewDiffTaskWithPrefix(renderRaw bool, cmd *exec.Cmd, prefix string) UpdateTask {
	if renderRaw {
		return NewRunCommandTaskWithPrefix(cmd, prefix)
	}
	return NewRunPtyTaskWithPrefix(cmd, prefix)
}
