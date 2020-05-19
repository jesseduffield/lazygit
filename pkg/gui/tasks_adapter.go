package gui

import (
	"os/exec"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

func (gui *Gui) newCmdTask(viewName string, cmd *exec.Cmd) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil // swallowing for now
	}

	_, height := view.Size()
	_, oy := view.Origin()

	manager := gui.getManager(view)

	r, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := manager.NewTask(manager.NewCmdTask(r, cmd, height+oy+10, nil)); err != nil {
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
		gui.renderString(gui.g, viewName, str)
		return nil
	}

	if err := manager.NewTask(f); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithoutScroll(viewName string, str string) error {
	view, err := gui.g.View(viewName)
	if err != nil {
		return nil // swallowing for now
	}

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
				view.Clear()
			},
			func() {
				gui.g.Update(func(*gocui.Gui) error {
					return nil
				})
			})
		gui.viewBufferManagerMap[view.Name()] = manager
	}

	return manager
}
