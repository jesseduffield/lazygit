package gui

import (
	"io"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/sirupsen/logrus"
)

func (gui *Gui) newCmdTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	cmdStr := strings.Join(cmd.Args, " ")
	gui.c.Log.WithField(
		"command",
		cmdStr,
	).Debug("RunCommand")

	// Mark the view as loading synchronously (before the task's goroutine runs
	// and before the next layout pass) so the layout doesn't clamp the scroll
	// position to the not-yet-loaded content.
	gui.getManager(view).StartLoading()
	// Hold the scrollbar at the height the view has now (the previous render),
	// while it still shows that render: once the re-render swaps in its first
	// partial paint the displayed buffer is briefly short, and we don't want the
	// thumb to shrink and snap back as the rest loads.
	view.FreezeScrollbarHeight()

	// Snapshot the view width here, on the UI thread, so the task goroutine
	// doesn't read the view's live dimensions while it streams output.
	spec := renderSpec{view: view, cmd: cmd, width: view.InnerWidth()}

	return gui.newTaskForRender(spec, prefix, cmdStr, gui.plainRender)
}

// plainRender runs the command as it is, with its output going straight into
// a pipe.
func (gui *Gui) plainRender(spec renderSpec) (startRender, onCloseRender) {
	var r io.ReadCloser
	start := func() (tasks.Cmd, io.Reader) {
		// The view wraps to this width; apply it here, on the task's goroutine
		// once the previous task has stopped, so it doesn't race that task's
		// still-running writes (see View.SetContentWidth).
		spec.view.SetContentWidth(spec.width)

		execCmd, pipe := startCmdWithPipe(spec.cmd, gui.c.Log)
		r = pipe
		return execCmd, pipe
	}

	onClose := func() {
		if r != nil {
			r.Close()
			r = nil
		}
	}

	return start, onClose
}

// startCmdWithPipe starts cmd with its stdout and stderr going to a single
// pipe, and returns the command along with the pipe's read end, in the shape
// that NewCmdTask expects from its start func. It never returns a nil reader,
// because NewCmdTask's scanner panics on one: when the pipe can't be created
// the command isn't started at all, and an empty reader is returned so that
// the task shuts down cleanly with the error in the log.
func startCmdWithPipe(cmd *exec.Cmd, log *logrus.Entry) (tasks.Cmd, io.ReadCloser) {
	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err)
		return tasks.ExecCmd{Cmd: cmd}, io.NopCloser(strings.NewReader(""))
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		log.Error(err)
	}

	return tasks.ExecCmd{Cmd: cmd}, r
}

func (gui *Gui) newStringTask(view *gocui.View, str string) error {
	// using str so that if rendering the exact same thing we don't reset the origin
	return gui.newStringTaskWithKey(view, str, str)
}

func (gui *Gui) newStringTaskWithoutScroll(view *gocui.View, str string) error {
	manager := gui.getManager(view)
	// Whatever the view was going to be put back to belonged to a re-render of its
	// content; this is a message instead, so there is nothing to put back.
	manager.DropRestoreForNextTask()

	f := func(tasks.TaskOpts) error {
		return gui.g.OnUIThreadAndWaitBackground(func() {
			gui.c.SetViewContent(view, str)
			gui.updateDiffSelectionVisibility(view, true)
			gui.reApplySearch(view)
		})
	}

	if err := manager.NewTask(f, manager.GetTaskKey()); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithScroll(view *gocui.View, str string, originX int, originY int) error {
	manager := gui.getManager(view)
	// Whatever the view was going to be put back to belonged to a re-render of its
	// content; this is a message instead, so there is nothing to put back.
	manager.DropRestoreForNextTask()

	f := func(tasks.TaskOpts) error {
		return gui.g.OnUIThreadAndWaitBackground(func() {
			gui.c.SetViewContent(view, str)
			view.SetOrigin(originX, originY)
			gui.updateDiffSelectionVisibility(view, true)
			gui.reApplySearch(view)
		})
	}

	if err := manager.NewTask(f, manager.GetTaskKey()); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newStringTaskWithKey(view *gocui.View, str string, key string) error {
	manager := gui.getManager(view)
	// Whatever the view was going to be put back to belonged to a re-render of its
	// content; this is a message instead, so there is nothing to put back.
	manager.DropRestoreForNextTask()

	f := func(tasks.TaskOpts) error {
		return gui.g.OnUIThreadAndWaitBackground(func() {
			gui.c.ResetViewOrigin(view)
			gui.c.SetViewContent(view, str)
			gui.updateDiffSelectionVisibility(view, true)
			gui.reApplySearch(view)
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
				// content-only render is enough — it skips the layout pass and redraws
				// only the cells that differ — and it's much cheaper than a full
				// layout-and-redraw on every read, which matters a lot when reading a
				// long diff, where reads happen repeatedly as the user scrolls.
				//
				// What this draws is more of the content than the pane held a moment
				// ago, so it is also where what is drawn over that content is worked
				// out again. The screenful the first paint reveals may not be enough
				// to say whether there is anything to select, and for a diff that
				// opens with a long diffstat it isn't.
				gui.c.OnUIThreadContentOnly(func() error {
					gui.updateDiffSelectionVisibility(view, false)
					return nil
				})
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

				gui.updateDiffSelectionVisibility(view, true)
				gui.clampDiffSelectionToContent(view)
				gui.reApplySearch(view)
			},
			func() {
				view.SetOrigin(0, 0)
			},
			view.BeginOffscreenRender,
			func() {
				view.SwapInOffscreenRender()

				// The content the pane is being given is on display from here on, so
				// what is drawn over it is settled against that content rather than
				// against the render before it.
				gui.updateDiffSelectionVisibility(view, false)
			},
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
