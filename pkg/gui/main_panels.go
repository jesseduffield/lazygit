package gui

import (
	"os/exec"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type viewUpdateOpts struct {
	title string

	// awkwardly calling this noWrap because of how hard Go makes it to have
	// a boolean option that defaults to true
	noWrap bool

	highlight bool

	task updateTask

	context types.Context
}

type refreshMainOpts struct {
	main      *viewUpdateOpts
	secondary *viewUpdateOpts
}

type updateTask interface {
	IsUpdateTask()
}

type renderStringTask struct {
	str string
}

func (t *renderStringTask) IsUpdateTask() {}

func NewRenderStringTask(str string) *renderStringTask {
	return &renderStringTask{str: str}
}

type renderStringWithoutScrollTask struct {
	str string
}

func (t *renderStringWithoutScrollTask) IsUpdateTask() {}

func NewRenderStringWithoutScrollTask(str string) *renderStringWithoutScrollTask {
	return &renderStringWithoutScrollTask{str: str}
}

type runCommandTask struct {
	cmd    *exec.Cmd
	prefix string
}

func (t *runCommandTask) IsUpdateTask() {}

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

func (t *runPtyTask) IsUpdateTask() {}

func NewRunPtyTask(cmd *exec.Cmd) *runPtyTask {
	return &runPtyTask{cmd: cmd}
}

// currently unused
// func (gui *Gui) createRunPtyTaskWithPrefix(cmd *exec.Cmd, prefix string) *runPtyTask {
// 	return &runPtyTask{cmd: cmd, prefix: prefix}
// }

func (gui *Gui) runTaskForView(view *gocui.View, task updateTask) error {
	switch v := task.(type) {
	case *renderStringTask:
		return gui.newStringTask(view, v.str)

	case *renderStringWithoutScrollTask:
		return gui.newStringTaskWithoutScroll(view, v.str)

	case *runCommandTask:
		return gui.newCmdTask(view, v.cmd, v.prefix)

	case *runPtyTask:
		return gui.newPtyTask(view, v.cmd, v.prefix)
	}

	return nil
}

func (gui *Gui) refreshMainView(opts *viewUpdateOpts, view *gocui.View) error {
	view.Title = opts.title
	view.Wrap = !opts.noWrap
	view.Highlight = opts.highlight
	context := opts.context
	if context == nil {
		context = gui.State.Contexts.Normal
	}
	gui.ViewContextMapSet(view.Name(), context)

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
