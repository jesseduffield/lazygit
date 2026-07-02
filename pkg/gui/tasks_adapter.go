package gui

import (
	"io"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

func (gui *Gui) newCmdTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	cmdStr := strings.Join(cmd.Args, " ")
	gui.c.Log.WithField(
		"command",
		cmdStr,
	).Debug("RunCommand")

	manager := gui.getManager(view)
	// Mark the view as loading synchronously (before the task's goroutine runs
	// and before the next layout pass) so the layout doesn't clamp the scroll
	// position to the not-yet-loaded content.
	manager.StartLoading()
	// Hold the scrollbar at the height the view has now (the previous render),
	// while it still shows that render: once the re-render swaps in its first
	// partial paint the displayed buffer is briefly short, and we don't want the
	// thumb to shrink and snap back as the rest loads.
	view.FreezeScrollbarHeight()

	var r io.ReadCloser
	start := func() (*exec.Cmd, io.Reader) {
		var err error
		r, err = cmd.StdoutPipe()
		if err != nil {
			gui.c.Log.Error(err)
			r = nil
		}
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			gui.c.Log.Error(err)
		}

		return cmd, r
	}

	onClose := func() {
		if r != nil {
			r.Close()
			r = nil
		}
	}

	linesToRead := gui.linesToReadFromCmdTask(view)
	// If a restore is pending for this content (returning to a focused main view
	// on escape), let the task re-establish the scroll position and selection as
	// it first paints. It also reads to end of input so a deep target line is
	// found and the scrollbar ends up accurate. See RenderRestore.
	restore := manager.GetRestoreForNextTask()
	if restore != nil {
		linesToRead.Restore = restore
		linesToRead.Total = -1
	}
	// New content (a different command than the view last showed) scrolls back to
	// the top, at the first paint; same content keeps its scroll, and a restore
	// places the scroll itself. See LinesToRead.ResetOrigin.
	linesToRead.ResetOrigin = restore == nil && cmdStr != manager.GetTaskKey()
	if err := manager.NewTask(manager.NewCmdTask(start, prefix, linesToRead, onClose), cmdStr); err != nil {
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

	if err := manager.NewTask(f, manager.GetTaskKey()); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithScroll(view *gocui.View, str string, originX int, originY int) error {
	manager := gui.getManager(view)

	f := func(tasks.TaskOpts) error {
		gui.c.SetViewContent(view, str)
		view.SetOrigin(originX, originY)
		return nil
	}

	if err := manager.NewTask(f, manager.GetTaskKey()); err != nil {
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
				// Called before showing the "loading..." indicator: clear the
				// displayed buffer so only "loading..." is shown. The actual content
				// is rendered off-screen (beginRender below) and swapped in, so it
				// never overwrites the displayed buffer incrementally.
				view.Reset()
			},
			func() {
				gui.render()
			},
			func() {
				// The content is fully loaded now, so let the scrollbar track it
				// directly again (it was held at the previous render's height while
				// loading, see FreezeScrollbarHeight).
				view.UnfreezeScrollbarHeight()

				// Need to check if the content of the view is well past the origin.
				linesHeight := view.ViewLinesHeight()
				_, originY := view.Origin()
				if linesHeight < originY {
					newOriginY := linesHeight

					view.SetOrigin(0, newOriginY)
				}

				view.FlushStaleCells()
			},
			func() {
				view.SetOrigin(0, 0)
			},
			view.BeginOffscreenRender,
			view.SwapInOffscreenRender,
			func() gocui.Task {
				return gui.c.GocuiGui().NewTask()
			},
		)
		gui.viewBufferManagerMap[view.Name()] = manager
	}

	return manager
}
