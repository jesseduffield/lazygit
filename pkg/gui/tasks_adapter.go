package gui

import (
	"os/exec"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/pty"
)

func (gui *Gui) newCmdTask(viewName string, cmd *exec.Cmd) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil // swallowing for now
	}

	_, height := view.Size()
	_, oy := view.Origin()

	manager := gui.getManager(view)

	if err := manager.NewTask(manager.NewCmdTask(cmd, height+oy+10)); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newPtyTask(viewName string, cmd *exec.Cmd) error {
	width, _ := gui.getMainView().Size()
	pager := gui.GitCommand.GetPager(width)

	if pager == "" {
		// if we're not using a custom pager we don't need to use a pty
		return gui.newCmdTask(viewName, cmd)
	}

	cmd.Env = append(cmd.Env, "GIT_PAGER="+pager)

	view, err := gui.g.View(viewName)
	if err != nil {
		return nil // swallowing for now
	}

	_, height := view.Size()
	_, oy := view.Origin()

	manager := gui.getManager(view)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	gui.State.Ptmx = ptmx
	onClose := func() { gui.State.Ptmx = nil }

	if err := gui.onResize(); err != nil {
		return err
	}

	if err := manager.NewTask(manager.NewPtyTask(ptmx, cmd, height+oy+10, onClose)); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newTask(viewName string, f func(chan struct{}) error) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil // swallowing for now
	}

	manager := gui.getManager(view)

	if err := manager.NewTask(f); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTask(viewName string, str string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil // swallowing for now
	}

	manager := gui.getManager(view)

	f := func(stop chan struct{}) error {
		return gui.renderString(gui.g, viewName, str)
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
				view.Clear()
			},
			func() {
				gui.g.Update(func(*gocui.Gui) error {
					gui.Log.Warn("updating view")
					return nil
				})
			})
		gui.viewBufferManagerMap[view.Name()] = manager
	}

	return manager
}
