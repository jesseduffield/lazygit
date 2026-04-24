package gui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
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

func (p ptyCmd) Wait() error             { return p.wait() }
func (p ptyCmd) String() string          { return p.cmd.String() }
func (p ptyCmd) GetProcess() *os.Process { return p.process }

// Some commands need to output for a terminal to active certain behaviour.
// For example,  git won't invoke the GIT_PAGER env var unless it thinks it's
// talking to a terminal. We typically write cmd outputs straight to a view,
// which is just an io.Reader. the pty package lets us wrap a command in a
// pseudo-terminal meaning we'll get the behaviour we want from the underlying
// command.
func (gui *Gui) newPtyTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	width := view.InnerWidth()

	// LAZYGIT_COLUMNS is a Windows-documented env var (see
	// docs/Custom_Pagers.md) for pager scripts that can't query the
	// terminal width directly. We set it on every platform so those
	// scripts remain portable.
	cmd.Env = append(cmd.Env, fmt.Sprintf("LAZYGIT_COLUMNS=%d", width))

	pager := gui.stateAccessor.GetPagerConfig().GetPagerCommand(width)
	externalDiffCommand := gui.stateAccessor.GetPagerConfig().GetExternalDiffCommand()
	useExtDiffGitConfig := gui.stateAccessor.GetPagerConfig().GetUseExternalDiffGitConfig()

	if pager == "" && externalDiffCommand == "" && !useExtDiffGitConfig {
		// If we're not using a custom pager nor external diff command, then we don't need to use a pty
		return gui.newCmdTask(view, cmd, prefix)
	}

	// Run the pty after layout so that it gets the correct size
	gui.afterLayout(func() error {
		// Need to get the width and the pager again because the layout might have
		// changed the size of the view
		width = view.InnerWidth()
		pager := gui.stateAccessor.GetPagerConfig().GetPagerCommand(width)

		cmdStr := strings.Join(cmd.Args, " ")

		// This communicates to pagers that we're in a very simple
		// terminal that they should not expect to have much capabilities.
		// Moving the cursor, clearing the screen, or querying for colors are among such "advanced" capabilities.
		// Context: https://github.com/jesseduffield/lazygit/issues/3419
		cmd.Env = removeExistingTermEnvVars(cmd.Env)
		cmd.Env = append(cmd.Env, "TERM=dumb")

		cmd.Env = append(cmd.Env, "GIT_PAGER="+pager)

		manager := gui.getManager(view)

		var p oscommands.Pty
		start := func() (tasks.Cmd, io.Reader) {
			cols, rows := gui.desiredPtySize(view)
			sp, err := oscommands.StartPty(cmd, cols, rows)
			if err != nil {
				gui.c.Log.Error(err)
				return tasks.ExecCmd{Cmd: cmd}, nil
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
			delete(gui.viewPtmxMap, view.Name())
			gui.Mutexes.PtyMutex.Unlock()
		}

		linesToRead := gui.linesToReadFromCmdTask(view)
		return manager.NewTask(manager.NewCmdTask(start, prefix, linesToRead, onClose), cmdStr)
	})

	return nil
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
