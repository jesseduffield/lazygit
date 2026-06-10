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

	// Snapshot the view width here, on the UI thread, so the task goroutine
	// doesn't read the view's live dimensions while it streams output. It's
	// applied inside start() below rather than now, because start() runs once
	// the previous task has stopped -- applying it here would race that task's
	// still-running writes (see View.SetContentWidth).
	contentWidth := view.InnerWidth()

	var r io.ReadCloser
	start := func() (tasks.Cmd, io.Reader) {
		view.SetContentWidth(contentWidth)

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

		return tasks.ExecCmd{Cmd: cmd}, r
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
		return gui.g.OnUIThreadAndWaitBackground(func() error {
			gui.c.SetViewContent(view, str)
			return nil
		})
	}

	if err := manager.NewTask(f, manager.GetTaskKey()); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithScroll(view *gocui.View, str string, originX int, originY int) error {
	manager := gui.getManager(view)

	f := func(tasks.TaskOpts) error {
		return gui.g.OnUIThreadAndWaitBackground(func() error {
			gui.c.SetViewContent(view, str)
			view.SetOrigin(originX, originY)
			return nil
		})
	}

	if err := manager.NewTask(f, manager.GetTaskKey()); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithKey(view *gocui.View, str string, key string) error {
	manager := gui.getManager(view)

	f := func(tasks.TaskOpts) error {
		return gui.g.OnUIThreadAndWaitBackground(func() error {
			gui.c.ResetViewOrigin(view)
			gui.c.SetViewContent(view, str)
			return nil
		})
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
				// As the task reads more lines, the only thing that changes is the
				// view's content (and its scrollbar); the window layout doesn't. So a
				// content-only render is enough, and it's much cheaper than a full
				// layout-and-redraw on every read - which matters a lot when reading
				// a long diff, where reads happen repeatedly as the user scrolls.
				gui.renderContentOnly()
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
				// A background task: rendering content into a view is display
				// work, not lazygit driving a git operation, so it must not
				// count towards being busy and block a repo switch. These
				// renders fire on nearly every focus/selection change, including
				// the context activation that happens right before a menu/prompt
				// handler runs (e.g. confirming worktree creation), which would
				// otherwise make the switch that handler triggers refuse itself.
				return gui.c.GocuiGui().NewBackgroundTask()
			},
			// Rendering is background work too (see above), so the view mutations
			// it bounces onto the UI thread mustn't count towards being busy.
			gui.g.OnUIThreadAndWaitBackground,
		)
		gui.viewBufferManagerMap[view.Name()] = manager
	}

	return manager
}
