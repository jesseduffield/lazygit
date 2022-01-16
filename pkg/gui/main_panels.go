package gui

import (
	"os/exec"

	"github.com/jesseduffield/gocui"
)

type viewUpdateOpts struct {
	title string

	// awkwardly calling this noWrap because of how hard Go makes it to have
	// a boolean option that defaults to true
	noWrap bool

	highlight bool

	task updateTask
}

type refreshMainOpts struct {
	main      *viewUpdateOpts
	secondary *viewUpdateOpts
}

// constants for updateTask's kind field
type TaskKind int

const (
	RENDER_STRING TaskKind = iota
	RENDER_STRING_WITHOUT_SCROLL
	RUN_COMMAND
	RUN_PTY
)

type updateTask interface {
	GetKind() TaskKind
}

type renderStringTask struct {
	str string
}

func (t *renderStringTask) GetKind() TaskKind {
	return RENDER_STRING
}

func NewRenderStringTask(str string) *renderStringTask {
	return &renderStringTask{str: str}
}

type renderStringWithoutScrollTask struct {
	str string
}

func (t *renderStringWithoutScrollTask) GetKind() TaskKind {
	return RENDER_STRING_WITHOUT_SCROLL
}

func NewRenderStringWithoutScrollTask(str string) *renderStringWithoutScrollTask {
	return &renderStringWithoutScrollTask{str: str}
}

type runCommandTask struct {
	cmd    *exec.Cmd
	prefix string
}

func (t *runCommandTask) GetKind() TaskKind {
	return RUN_COMMAND
}

func NewRunCommandTask(cmd *exec.Cmd) *runCommandTask {
	return &runCommandTask{cmd: cmd}
}

func NewRunCommandTaskWithPrefix(cmd *exec.Cmd, prefix string) *runCommandTask {
	return &runCommandTask{cmd: cmd, prefix: prefix}
}

type runPtyTask struct {
	cmd    *exec.Cmd
	prefix string
}

func (t *runPtyTask) GetKind() TaskKind {
	return RUN_PTY
}

func NewRunPtyTask(cmd *exec.Cmd) *runPtyTask {
	return &runPtyTask{cmd: cmd}
}

// currently unused
// func (gui *Gui) createRunPtyTaskWithPrefix(cmd *exec.Cmd, prefix string) *runPtyTask {
// 	return &runPtyTask{cmd: cmd, prefix: prefix}
// }

func (gui *Gui) runTaskForView(view *gocui.View, task updateTask) error {
	switch task.GetKind() {
	case RENDER_STRING:
		specificTask := task.(*renderStringTask)
		return gui.newStringTask(view, specificTask.str)

	case RENDER_STRING_WITHOUT_SCROLL:
		specificTask := task.(*renderStringWithoutScrollTask)
		return gui.newStringTaskWithoutScroll(view, specificTask.str)

	case RUN_COMMAND:
		specificTask := task.(*runCommandTask)
		return gui.newCmdTask(view, specificTask.cmd, specificTask.prefix)

	case RUN_PTY:
		specificTask := task.(*runPtyTask)
		return gui.newPtyTask(view, specificTask.cmd, specificTask.prefix)
	}

	return nil
}

func (gui *Gui) refreshMainView(opts *viewUpdateOpts, view *gocui.View) error {
	view.Title = opts.title
	view.Wrap = !opts.noWrap
	view.Highlight = opts.highlight

	if err := gui.runTaskForView(view, opts.task); err != nil {
		gui.c.Log.Error(err)
		return nil
	}

	return nil
}

func (gui *Gui) refreshMainViews(opts refreshMainOpts) error {
	if opts.main != nil {
		if err := gui.refreshMainView(opts.main, gui.Views.Main); err != nil {
			return err
		}
	}

	if opts.secondary != nil {
		if err := gui.refreshMainView(opts.secondary, gui.Views.Secondary); err != nil {
			return err
		}
	}

	gui.splitMainPanel(opts.secondary != nil)

	return nil
}

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel
}

func (gui *Gui) isMainPanelSplit() bool {
	return gui.State.SplitMainPanel
}
