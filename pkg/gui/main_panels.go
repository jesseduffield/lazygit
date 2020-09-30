package gui

import "os/exec"

type viewUpdateOpts struct {
	title string

	// awkwardly calling this noWrap because of how hard Go makes it to have
	// a boolean option that defaults to true
	noWrap bool

	highlight bool

	task updateTask
}

type coordinates struct {
	x int
	y int
}

type refreshMainOpts struct {
	main      *viewUpdateOpts
	secondary *viewUpdateOpts
}

// constants for updateTask's kind field
const (
	RENDER_STRING = iota
	RENDER_STRING_WITHOUT_SCROLL
	RUN_FUNCTION
	RUN_COMMAND
	RUN_PTY
)

type updateTask interface {
	GetKind() int
}

type renderStringTask struct {
	str string
}

func (t *renderStringTask) GetKind() int {
	return RENDER_STRING
}

func (gui *Gui) createRenderStringTask(str string) *renderStringTask {
	return &renderStringTask{str: str}
}

type renderStringWithoutScrollTask struct {
	str string
}

func (t *renderStringWithoutScrollTask) GetKind() int {
	return RENDER_STRING_WITHOUT_SCROLL
}

func (gui *Gui) createRenderStringWithoutScrollTask(str string) *renderStringWithoutScrollTask {
	return &renderStringWithoutScrollTask{str: str}
}

type runCommandTask struct {
	cmd    *exec.Cmd
	prefix string
}

func (t *runCommandTask) GetKind() int {
	return RUN_COMMAND
}

func (gui *Gui) createRunCommandTask(cmd *exec.Cmd) *runCommandTask {
	return &runCommandTask{cmd: cmd}
}

func (gui *Gui) createRunCommandTaskWithPrefix(cmd *exec.Cmd, prefix string) *runCommandTask {
	return &runCommandTask{cmd: cmd, prefix: prefix}
}

type runPtyTask struct {
	cmd    *exec.Cmd
	prefix string
}

func (t *runPtyTask) GetKind() int {
	return RUN_PTY
}

func (gui *Gui) createRunPtyTask(cmd *exec.Cmd) *runPtyTask {
	return &runPtyTask{cmd: cmd}
}

func (gui *Gui) createRunPtyTaskWithPrefix(cmd *exec.Cmd, prefix string) *runPtyTask {
	return &runPtyTask{cmd: cmd, prefix: prefix}
}

type runFunctionTask struct {
	f func(chan struct{}) error
}

func (t *runFunctionTask) GetKind() int {
	return RUN_FUNCTION
}

func (gui *Gui) createRunFunctionTask(f func(chan struct{}) error) *runFunctionTask {
	return &runFunctionTask{f: f}
}

func (gui *Gui) runTaskForView(viewName string, task updateTask) error {
	switch task.GetKind() {
	case RENDER_STRING:
		specificTask := task.(*renderStringTask)
		return gui.newStringTask(viewName, specificTask.str)

	case RENDER_STRING_WITHOUT_SCROLL:
		specificTask := task.(*renderStringWithoutScrollTask)
		return gui.newStringTaskWithoutScroll(viewName, specificTask.str)

	case RUN_FUNCTION:
		specificTask := task.(*runFunctionTask)
		return gui.newTask(viewName, specificTask.f)

	case RUN_COMMAND:
		specificTask := task.(*runCommandTask)
		return gui.newCmdTask(viewName, specificTask.cmd, specificTask.prefix)

	case RUN_PTY:
		specificTask := task.(*runPtyTask)
		return gui.newPtyTask(viewName, specificTask.cmd, specificTask.prefix)
	}

	return nil
}

func (gui *Gui) refreshMainView(opts *viewUpdateOpts, viewName string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		gui.Log.Error(err)
		return nil
	}

	view.Title = opts.title
	view.Wrap = !opts.noWrap
	view.Highlight = opts.highlight

	if err := gui.runTaskForView(viewName, opts.task); err != nil {
		gui.Log.Error(err)
		return nil
	}

	return nil
}

func (gui *Gui) refreshMainViews(opts refreshMainOpts) error {
	if opts.main != nil {
		if err := gui.refreshMainView(opts.main, "main"); err != nil {
			return err
		}
	}

	gui.splitMainPanel(opts.secondary != nil)

	if opts.secondary != nil {
		if err := gui.refreshMainView(opts.secondary, "secondary"); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel

	// no need to set view on bottom when splitMainPanel is false: it will have zero size anyway thanks to our view arrangement code.
	if splitMainPanel {
		_, _ = gui.g.SetViewOnTop("secondary")
	}
}

func (gui *Gui) isMainPanelSplit() bool {
	return gui.State.SplitMainPanel
}
