package gui

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// withPtyGitConfig returns args with extra git configuration for commands
// that render into a pty. On Windows, such a command is terminated at an
// arbitrary point of its execution when its task stops: tearing down the
// pseudoconsole delivers CTRL_CLOSE_EVENT, which git leaves to the default
// handler, which just calls ExitProcess. git's automatic index refresh
// (diff.autoRefreshIndex, on by default) takes index.lock at the end of a
// diff against the worktree to write back refreshed stat information —
// GIT_OPTIONAL_LOCKS does not cover this lock — and a termination landing
// in that window leaves a stale index.lock behind that the next git command
// chokes on. So don't let pty-rendered commands refresh the index;
// lazygit's foreground `git status` refreshes, which never run in a pty,
// keep the stat cache fresh instead.
//
// On Unix a stopped pty child gets SIGTERM, and git's signal handlers remove
// its lock files, so the refresh can stay enabled there and keep healing
// stale stat info.
func withPtyGitConfig(args []string, goos string) []string {
	if goos != "windows" {
		return args
	}
	// Most pty commands are direct git invocations, but the user-configured
	// ones can be arbitrary command lines (e.g. a branchLogCmd wrapping git
	// in `sh -c`), and injecting git flags into those would corrupt them.
	// Only direct git invocations get the config; that loses nothing, since
	// the wrapped commands are log commands, which never take the index
	// lock. (For direct invocations other than worktree diffs the config is
	// simply a no-op.)
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(args[0])), ".exe")
	if base != "git" {
		return args
	}
	result := make([]string, 0, len(args)+2)
	result = append(result, args[0])
	result = append(result, "-c", "diff.autoRefreshIndex=false")
	return append(result, args[1:]...)
}
