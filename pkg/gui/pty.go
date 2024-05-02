//go:build !windows
// +build !windows

package gui

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

func (gui *Gui) desiredPtySize() *pty.Winsize {
	width, height := gui.Views.Main.Size()

	return &pty.Winsize{Cols: uint16(width), Rows: uint16(height)}
}

func (gui *Gui) onResize() error {
	gui.Mutexes.PtyMutex.Lock()
	defer gui.Mutexes.PtyMutex.Unlock()

	for _, ptmx := range gui.viewPtmxMap {
		// TODO: handle resizing properly: we need to actually clear the main view
		// and re-read the output from our pty. Or we could just re-run the original
		// command from scratch
		if err := pty.Setsize(ptmx, gui.desiredPtySize()); err != nil {
			return utils.WrapError(err)
		}
	}

	return nil
}

// Some commands need to output for a terminal to active certain behaviour.
// For example,  git won't invoke the GIT_PAGER env var unless it thinks it's
// talking to a terminal. We typically write cmd outputs straight to a view,
// which is just an io.Reader. the pty package lets us wrap a command in a
// pseudo-terminal meaning we'll get the behaviour we want from the underlying
// command.
func (gui *Gui) newPtyTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	width, _ := gui.Views.Main.Size()
	pager := gui.git.Config.GetPager(width)
	externalDiffCommand := gui.Config.GetUserConfig().Git.Paging.ExternalDiffCommand

	if pager == "" && externalDiffCommand == "" {
		// if we're not using a custom pager we don't need to use a pty
		return gui.newCmdTask(view, cmd, prefix)
	}

	cmdStr := strings.Join(cmd.Args, " ")

	// This communicates to pagers that we're in a very simple
	// terminal that they should not expect to have much capabilities.
	// Moving the cursor, clearing the screen, or querying for colors are among such "advanced" capabilities.
	// Context: https://github.com/jesseduffield/lazygit/issues/3419
	cmd.Env = removeExistingTermEnvVars(cmd.Env)
	cmd.Env = append(cmd.Env, "TERM=dumb")

	cmd.Env = append(cmd.Env, "GIT_PAGER="+pager)

	manager := gui.getManager(view)

	var ptmx *os.File
	start := func() (*exec.Cmd, io.Reader) {
		var err error
		ptmx, err = pty.StartWithSize(cmd, gui.desiredPtySize())
		if err != nil {
			gui.c.Log.Error(err)
		}

		gui.Mutexes.PtyMutex.Lock()
		gui.viewPtmxMap[view.Name()] = ptmx
		gui.Mutexes.PtyMutex.Unlock()

		return cmd, ptmx
	}

	onClose := func() {
		gui.Mutexes.PtyMutex.Lock()
		ptmx.Close()
		delete(gui.viewPtmxMap, view.Name())
		gui.Mutexes.PtyMutex.Unlock()
	}

	linesToRead := gui.linesToReadFromCmdTask(view)
	if err := manager.NewTask(manager.NewCmdTask(start, prefix, linesToRead, onClose), cmdStr); err != nil {
		return err
	}

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
