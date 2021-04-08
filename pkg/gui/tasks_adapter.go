package gui

import (
	"os/exec"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

func (gui *Gui) newCmdTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	gui.Log.WithField(
		"command",
		strings.Join(cmd.Args, " "),
	).Debug("RunCommand")

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
