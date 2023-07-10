package gui

import (
	"io"
	"os/exec"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

func (gui *Gui) newCmdTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	cmdStr := strings.Join(cmd.Args, " ")
	gui.c.Log.WithField(
		"command",
		cmdStr,
	).Debug("RunCommand")

	manager := gui.getManager(view)

	start := func() (*exec.Cmd, io.Reader) {
		r, err := cmd.StdoutPipe()
		if err != nil {
			gui.c.Log.Error(err)
		}
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			gui.c.Log.Error(err)
		}

		return cmd, r
	}

	linesToRead := gui.linesToReadFromCmdTask(view)
	if err := manager.NewTask(manager.NewCmdTask(start, prefix, linesToRead, nil), cmdStr); err != nil {
		gui.c.Log.Error(err)
	}

	return nil
}

func (gui *Gui) newStringTask(view *gocui.View, str string) error {
	// using str so that if rendering the exact same thing we don't reset the origin
	return gui.newStringTaskWithKey(view, str, str)
}

func (gui *Gui) newStringTaskWithoutScroll(view *gocui.View, str string) error {
	manager := gui.getManager(view)

	f := func(tasks.TaskOpts) error {
		gui.c.SetViewContent(view, str)
		return nil
	}

	// Using empty key so that on subsequent calls we won't reset the view's origin.
	// Note this means that we will be scrolling back to the top if we're switching from a different key
	if err := manager.NewTask(f, ""); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithScroll(view *gocui.View, str string, originX int, originY int) error {
	manager := gui.getManager(view)

	f := func(tasks.TaskOpts) error {
		gui.c.SetViewContent(view, str)
		_ = view.SetOrigin(originX, originY)
		return nil
	}

	if err := manager.NewTask(f, ""); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithKey(view *gocui.View, str string, key string) error {
	manager := gui.getManager(view)

	f := func(tasks.TaskOpts) error {
		gui.c.ResetViewOrigin(view)
		gui.c.SetViewContent(view, str)
		return nil
	}

	if err := manager.NewTask(f, key); err != nil {
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
				gui.render()
			},
			func() {
				// Need to check if the content of the view is well past the origin.
				linesHeight := view.ViewLinesHeight()
				_, originY := view.Origin()
				if linesHeight < originY {
					newOriginY := linesHeight

					err := view.SetOrigin(0, newOriginY)
					if err != nil {
						panic(err)
					}
				}

				view.FlushStaleCells()
			},
			func() {
				_ = view.SetOrigin(0, 0)
			},
			func() gocui.Task {
				return gui.c.GocuiGui().NewTask()
			},
		)
		gui.viewBufferManagerMap[view.Name()] = manager
	}

	return manager
}
