package gui

import (
	"io"
	"os"
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) desiredPtySize(view *gocui.View) (cols, rows uint16) {
	width, height := view.InnerSize()
	return uint16(width), uint16(height)
}

func (gui *Gui) onResize() error {
	gui.Mutexes.PtyMutex.Lock()
	defer gui.Mutexes.PtyMutex.Unlock()

	for viewName, p := range gui.viewPtmxMap {
		// TODO: handle resizing properly: we need to actually clear the main view
		// and re-read the output from our pty. Or we could just re-run the original
		// command from scratch
		view, _ := gui.g.View(viewName)
		cols, rows := gui.desiredPtySize(view)
		if err := p.Resize(cols, rows); err != nil {
			return utils.WrapError(err)
		}
	}

	return nil
}

// ptyCmd adapts an oscommands.StartedPty result into the tasks.Cmd shape.
// On Windows the original *exec.Cmd was never Start()ed, so we go through
// the explicit Process handle rather than cmd.Process.
type ptyCmd struct {
	cmd     *exec.Cmd
	process *os.Process
	wait    func() error
}

func (p ptyCmd) Wait() error      { return p.wait() }
func (p ptyCmd) String() string   { return p.cmd.String() }
func (p ptyCmd) Terminate() error { return oscommands.TerminateProcessGracefully(p.process) }

// ptyRender runs the command in a pseudo-terminal. git invokes the stdin filter
// named by GIT_PAGER only when it talks to a terminal, and a renderer reads the
// width it lays out to off it.
//
// Must be called on the UI thread: it reads the view's dimensions, which the
// layout writes.
func (gui *Gui) ptyRender(spec renderSpec) (startRender, onCloseRender) {
	view := spec.view
	cmd := spec.cmd

	// git runs the stdin filter itself, as the pager it is told about here.
	// Named even when there is none, so that git doesn't reach for the user's
	// core.pager instead.
	cmd.Env = append(cmd.Env, "GIT_PAGER="+spec.stdinFilter)

	cols, rows := gui.desiredPtySize(view)

	var p oscommands.Pty
	var fallbackPipe io.ReadCloser
	start := func() (tasks.Cmd, io.Reader) {
		// The pty (and diff renderer) wrap to this width; apply it here, on the
		// task's goroutine once the previous task has stopped, so it doesn't
		// race that task's writes (see View.SetContentWidth).
		view.SetContentWidth(spec.width)

		sp, err := oscommands.StartPty(cmd, cols, rows)
		if err != nil {
			gui.c.Log.Error(err)
			// Fall back to running the command without a pty: the diff renderer is
			// lost, but the command's output still renders.
			execCmd, pipe := startCmdWithPipe(cmd, gui.c.Log)
			fallbackPipe = pipe
			return execCmd, pipe
		}
		p = sp.Pty

		gui.Mutexes.PtyMutex.Lock()
		gui.viewPtmxMap[view.Name()] = p
		gui.Mutexes.PtyMutex.Unlock()

		return ptyCmd{cmd: cmd, process: sp.Process, wait: sp.Wait}, p
	}

	onClose := func() {
		gui.Mutexes.PtyMutex.Lock()
		if p != nil {
			p.Close()
		}
		if fallbackPipe != nil {
			fallbackPipe.Close()
			fallbackPipe = nil
		}
		delete(gui.viewPtmxMap, view.Name())
		gui.Mutexes.PtyMutex.Unlock()
	}

	return start, onClose
}
