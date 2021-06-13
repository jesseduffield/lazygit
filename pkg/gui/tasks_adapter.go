package gui

import (
	"github.com/jesseduffield/gocui"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

func (gui *Gui) newCmdTask(view *gocui.View, cmdObj ICmdObj, prefix string) error {
	gui.Log.WithField("command", cmdObj.ToString()).Debug("RunCommand")

	_, height := view.Size()
	_, oy := view.Origin()

	manager := gui.getManager(view)

	cmd := cmdObj.GetCmd()
	r, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := manager.NewTask(manager.NewCmdTask(r, cmd, prefix, height+oy+10, nil)); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newTask(view *gocui.View, f func(chan struct{}) error) error {
	manager := gui.getManager(view)

	if err := manager.NewTask(f); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTask(view *gocui.View, str string) error {
	manager := gui.getManager(view)

	f := func(stop chan struct{}) error {
		gui.renderString(view, str)
		return nil
	}

	if err := manager.NewTask(f); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithoutScroll(view *gocui.View, str string) error {
	manager := gui.getManager(view)

	f := func(stop chan struct{}) error {
		gui.setViewContent(view, str)
		return nil
	}

	if err := manager.NewTask(f); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) getManager(view *gocui.View) *tasks.ViewBufferManager {
	manager, ok := gui.viewBufferManagerMap[view.Name()]
	if !ok {
		manager = tasks.NewViewBufferManager(
			gui.Log,
			view,
			func() {
				// we could clear here, but that actually has the effect of causing a flicker
				// where the view may contain no content momentarily as the gui refreshes.
				// Instead, we're rewinding the write pointer so that we will just start
				// overwriting the existing content from the top down. Once we've reached
				// the end of the content do display, we call view.FlushStaleCells() to
				// clear out the remaining content from the previous render.
				view.Reset()
			},
			func() {
				gui.g.Update(func(*gocui.Gui) error {
					return nil
				})
			},
			func() {
				view.FlushStaleCells()
			},
		)
		gui.viewBufferManagerMap[view.Name()] = manager
	}

	return manager
}
