package gui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/samber/lo"
)

// renderSpec is what setting up a render needs to know about it, whichever way
// the diff renderer is given the command's output.
type renderSpec struct {
	view *gocui.View
	cmd  *exec.Cmd
	// The width the renderer lays its rendering out to, and the width the view
	// counts its own wrapping against. Read after the layout pass, since that
	// is what settles it.
	width int
	// The configured stdin filter, empty unless one is configured. An external
	// diff renderer is named to git in the environment instead, so it doesn't
	// appear here; nothing else is a command of its own.
	stdinFilter string
}

// newRenderTask renders cmd's output into view, through the diff renderer the
// user has configured.
//
// Some commands need to output for a terminal to active certain behaviour.
// For example, git won't invoke the GIT_PAGER env var unless it thinks it's
// talking to a terminal. We typically write cmd outputs straight to a view,
// which is just an io.Reader. the pty package lets us wrap a command in a
// pseudo-terminal meaning we'll get the behaviour we want from the underlying
// command.
func (gui *Gui) newRenderTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	width := view.InnerWidth()

	// Set LAZYGIT_COLUMNS for diff renderer scripts that can't query the terminal width directly.
	cmd.Env = append(cmd.Env, fmt.Sprintf("LAZYGIT_COLUMNS=%d", width))

	if gui.stateAccessor.GetDiffRendererConfigManager().GetDiffRendererType() == config.DiffRendererType_RawGit {
		// If we're not using a custom diff renderer, then we don't need to use a pty
		return gui.newCmdTask(view, cmd, prefix)
	}

	cmd.Args = withPtyGitConfig(cmd.Args, runtime.GOOS)

	// Mark the view as loading synchronously now, before the layout pass: the
	// actual task is created in afterLayout (below), which runs after layout, so
	// without this the next layout pass would clamp the scroll position to the
	// not-yet-loaded content.
	gui.getManager(view).StartLoading()
	// Hold the scrollbar at its current height while the re-render loads, so the
	// thumb doesn't shrink and snap back when the first partial paint swaps in
	// (see the matching call in newCmdTask).
	view.FreezeScrollbarHeight()

	// Run the render after layout so that it gets the correct size
	gui.afterLayout(func() error {
		// Need to get the width and the renderer command again because the layout
		// might have changed the size of the view
		width = view.InnerWidth()
		diffRendererConfigManager := gui.stateAccessor.GetDiffRendererConfigManager()
		stdinFilter := diffRendererConfigManager.GetStdinFilterCommand(width)
		externalDiff := diffRendererConfigManager.GetExternalDiffCommand(gui.c.UserConfig().Git.DiffContextSize, width)

		cmdStr := strings.Join(cmd.Args, " ")

		// This communicates to diff renderers that we're in a very simple
		// terminal that they should not expect to have much capabilities.
		// Moving the cursor, clearing the screen, or querying for colors are among such "advanced" capabilities.
		// Context: https://github.com/jesseduffield/lazygit/issues/3419
		cmd.Env = removeExistingTermEnvVars(cmd.Env)
		cmd.Env = append(cmd.Env, "TERM=dumb")

		cmd.Env = append(cmd.Env, "GIT_PAGER="+stdinFilter)

		// An external diff command reaches git the same way, rather than as a
		// diff.external config on the command: it is only here, after the
		// layout, that the width it is to render at is known. An empty one
		// means the user wants git's own diff.external config to apply, so
		// leave the variable unset in that case; git takes it being set at all
		// as an instruction, however little it says.
		if externalDiff != "" {
			cmd.Env = append(cmd.Env, "GIT_EXTERNAL_DIFF="+externalDiff)
		}

		manager := gui.getManager(view)

		spec := renderSpec{
			view:        view,
			cmd:         cmd,
			width:       width,
			stdinFilter: stdinFilter,
		}
		var start startRender
		var onClose onCloseRender
		if rendersThroughAPipe() {
			start, onClose = gui.pipedRender(spec)
		} else {
			start, onClose = gui.ptyRender(spec)
		}

		linesToRead := gui.linesToReadFromCmdTask(view)
		return manager.NewTask(manager.NewCmdTask(start, prefix, linesToRead, onClose), cmdStr)
	})

	return nil
}

// The start and onClose functions a render hands to its task: how to get the
// command running and the output reader for it, and how to tear it down again
// once the task is stopped.
type (
	startRender   func() (tasks.Cmd, io.Reader)
	onCloseRender func()
)

// renderWithoutPtyEnvVar makes a render take the piped path on a platform that
// would otherwise use a pty, so that tests can exercise it anywhere.
const renderWithoutPtyEnvVar = "LAZYGIT_RENDER_WITHOUT_PTY"

// rendersThroughAPipe reports whether a render feeds the diff renderer the
// command's output through a pipe rather than running it in a pty.
//
// On Windows it has to. ConPTY doesn't pass a command's output through; it
// parses it into a screen buffer and re-encodes that for the terminal side,
// and it hands a sequence it can't represent there the moment it parses it,
// separately from the text around it. So what a renderer writes is not what
// lazygit reads. A pipe carries the bytes as the renderer wrote them.
//
// Everywhere else the pty is kept, since a renderer can read the width it
// should lay out to off it, and a configuration that doesn't name a width would
// otherwise render at whatever width the renderer falls back to.
func rendersThroughAPipe() bool {
	return runtime.GOOS == "windows" || os.Getenv(renderWithoutPtyEnvVar) != ""
}

// pipedRender feeds the diff renderer the command's output through a pipe.
//
// A stdin filter has to become a command of our own here: git only invokes the
// one named by GIT_PAGER when it thinks it is talking to a terminal, so with a
// pipe the filter would never run. An external diff renderer is git's own
// business — it is named in the environment and git runs it per file — so
// there the command still runs alone.
//
// Must be called on the UI thread; see ptyRender.
func (gui *Gui) pipedRender(spec renderSpec) (startRender, onCloseRender) {
	view := spec.view
	cmd := spec.cmd

	var pipe io.ReadCloser
	start := func() (tasks.Cmd, io.Reader) {
		// See the matching call in ptyRender for why this happens here.
		view.SetContentWidth(spec.width)

		tasks.DumpStreamNote("view %s: no pty, content width=%d, renderer: %s",
			view.Name(), spec.width, spec.stdinFilter)

		if spec.stdinFilter == "" {
			execCmd, reader := startCmdWithPipe(cmd, gui.c.Log)
			pipe = reader
			return execCmd, reader
		}

		// The filter runs in a plain shell, without lazygit's shell functions
		// sourced, since that is the shell git would have run it in. It is
		// handed git's environment for the same reason: as git's child it
		// would have inherited exactly that.
		pipeline, reader, err := gui.os.StartPipeline(
			gui.os.Cmd.NewFromCmd(cmd).DontLog(),
			gui.os.Cmd.NewShell(spec.stdinFilter, "").SetEnviron(cmd.Env).DontLog(),
		)
		if err != nil {
			gui.c.Log.Error(err)
			// The command has been started and stopped again by now, so it
			// can't be run a second time without the renderer. Show what went
			// wrong where the diff would have been.
			return tasks.ExecCmd{Cmd: cmd}, strings.NewReader(err.Error())
		}
		pipe = reader
		return pipeline, reader
	}

	onClose := func() {
		// Closing the reader is what brings the pipeline down: the renderer's
		// next write fails, so it exits, and git's write into the pipe the
		// renderer was reading fails in turn.
		if pipe != nil {
			pipe.Close()
			pipe = nil
		}
	}

	return start, onClose
}

func removeExistingTermEnvVars(env []string) []string {
	return lo.Filter(env, func(envVar string, _ int) bool {
		return !isTermEnvVar(envVar)
	})
}

// Terminals set a variety of different environment variables
// to identify themselves to processes. This list should catch the most common among them.
func isTermEnvVar(envVar string) bool {
	return strings.HasPrefix(envVar, "TERM=") ||
		strings.HasPrefix(envVar, "TERM_PROGRAM=") ||
		strings.HasPrefix(envVar, "TERM_PROGRAM_VERSION=") ||
		strings.HasPrefix(envVar, "TERMINAL_EMULATOR=") ||
		strings.HasPrefix(envVar, "TERMINAL_NAME=") ||
		strings.HasPrefix(envVar, "TERMINAL_VERSION_")
}
