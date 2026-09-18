package gui

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/samber/lo"
)

// renderSpec describes a render of a command's output into a view: what a way
// of running the command (see runRender) needs to know about it.
type renderSpec struct {
	view *gocui.View
	cmd  *exec.Cmd
	// The width the renderer lays its rendering out to, and the width the view
	// counts its own wrapping against. Read after the layout pass, which
	// settles it.
	width int
	// The configured stdin filter, empty unless one is configured. How git gets
	// to run it depends on the way the render runs. An external diff renderer
	// is named to git in the environment before the render is set up, whichever
	// way it runs, so it doesn't appear here; nothing else is a command of its
	// own.
	stdinFilter string
}

// newRenderTask renders cmd's output into view, through the diff renderer the
// user has configured. The renderer lays its rendering out to the width of the
// view, which only the layout settles, so the task is created after it.
func (gui *Gui) newRenderTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
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
		// The layout may have changed the size of the view, so only now is the
		// width to render at known, and with it the renderer command.
		width := view.InnerWidth()
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

		// An external diff command is named to git here, in the environment,
		// because the width it renders at is only known after the layout, and
		// the command's arguments were settled before it. An empty command
		// means the user wants git's own diff.external config to apply, so
		// leave the variable unset in that case; git takes it being set at all
		// as an instruction, however little it says.
		if externalDiff != "" {
			cmd.Env = append(cmd.Env, "GIT_EXTERNAL_DIFF="+externalDiff)
		}

		spec := renderSpec{
			view:        view,
			cmd:         cmd,
			width:       width,
			stdinFilter: stdinFilter,
		}
		return gui.newTaskForRender(spec, prefix, cmdStr, gui.ptyRender)
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

// runRender is a way of running a render's command and getting at its output:
// plainly, or in a pty. It returns the functions the task drives the command by.
type runRender func(spec renderSpec) (startRender, onCloseRender)

// newTaskForRender creates the task that reads the render's output into its
// view, running the command the given way. key names what is rendered, so that
// a re-render of the same content can be told from a render of other content.
func (gui *Gui) newTaskForRender(spec renderSpec, prefix string, key string, run runRender) error {
	setColumnsEnvVar(spec.cmd, spec.width)

	start, onClose := run(spec)

	manager := gui.getManager(spec.view)
	linesToRead := gui.linesToReadFromCmdTask(spec.view)
	return manager.NewTask(manager.NewCmdTask(start, prefix, linesToRead, onClose), key)
}

// setColumnsEnvVar tells a command how wide the view its output goes into is.
// git reads COLUMNS in preference to the size of the terminal it is talking to,
// and lays the diffstat graph out to it; a diff renderer with no terminal to
// ask may read it too (difftastic and diff-so-fancy do, delta does not). A
// command told nothing renders for 80 columns.
func setColumnsEnvVar(cmd *exec.Cmd, width int) {
	cmd.Env = append(cmd.Env, fmt.Sprintf("COLUMNS=%d", width))
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
