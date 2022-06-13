package gui

import (
	"os/exec"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type viewUpdateOpts struct {
	title string

	task updateTask
}

type refreshMainOpts struct {
	pair      MainContextPair
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

func (gui *Gui) moveMainContextPairToTop(pair MainContextPair) {
	gui.setWindowContext(pair.main)
	gui.moveToTopOfWindow(pair.main)
	if pair.secondary != nil {
		gui.setWindowContext(pair.secondary)
		gui.moveToTopOfWindow(pair.secondary)
	}
}

func (gui *Gui) refreshMainView(opts *viewUpdateOpts, context types.Context) error {
	view := context.GetView()

	if opts.title != "" {
		view.Title = opts.title
	}

	if err := gui.runTaskForView(view, opts.task); err != nil {
		gui.c.Log.Error(err)
		return nil
	}

	return nil
}

type MainContextPair struct {
	main      types.Context
	secondary types.Context
}

func (gui *Gui) normalMainContextPair() MainContextPair {
	return MainContextPair{
		main:      gui.State.Contexts.Normal,
		secondary: gui.State.Contexts.NormalSecondary,
	}
}

func (gui *Gui) stagingMainContextPair() MainContextPair {
	return MainContextPair{
		main:      gui.State.Contexts.Staging,
		secondary: gui.State.Contexts.StagingSecondary,
	}
}

func (gui *Gui) patchBuildingMainContextPair() MainContextPair {
	return MainContextPair{
		main:      gui.State.Contexts.CustomPatchBuilder,
		secondary: gui.State.Contexts.CustomPatchBuilderSecondary,
	}
}

func (gui *Gui) mergingMainContextPair() MainContextPair {
	return MainContextPair{
		main:      gui.State.Contexts.Merging,
		secondary: nil,
	}
}

func (gui *Gui) allMainContextPairs() []MainContextPair {
	return []MainContextPair{
		gui.normalMainContextPair(),
		gui.stagingMainContextPair(),
		gui.patchBuildingMainContextPair(),
		gui.mergingMainContextPair(),
	}
}

func (gui *Gui) refreshMainViews(opts refreshMainOpts) error {
	// need to reset scroll positions of all other main views
	for _, pair := range gui.allMainContextPairs() {
		if pair.main != opts.pair.main {
			_ = pair.main.GetView().SetOrigin(0, 0)
		}
		if pair.secondary != nil && pair.secondary != opts.pair.secondary {
			_ = pair.secondary.GetView().SetOrigin(0, 0)
		}
	}

	if opts.main != nil {
		if err := gui.refreshMainView(opts.main, opts.pair.main); err != nil {
			return err
		}
	}

	if opts.secondary != nil {
		if err := gui.refreshMainView(opts.secondary, opts.pair.secondary); err != nil {
			return err
		}
	} else if opts.pair.secondary != nil {
		opts.pair.secondary.GetView().Clear()
	}

	gui.moveMainContextPairToTop(opts.pair)

	gui.splitMainPanel(opts.secondary != nil)

	return nil
}

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel
}

func (gui *Gui) isMainPanelSplit() bool {
	return gui.State.SplitMainPanel
}
