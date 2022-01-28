package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Currently there is a bug where if we switch to a subprocess from within
// WithWaitingStatus we get stuck there and can't return to lazygit. We could
// fix this bug, or just stop running subprocesses from within there, given that
// we don't need to see a loading status if we're in a subprocess.
// TODO: work out if we actually need to use a shell command here
func (gui *Gui) withGpgHandling(cmdObj oscommands.ICmdObj, waitingStatus string, onSuccess func() error) error {
	gui.logCommand(cmdObj.ToString(), true)

	useSubprocess := gui.Git.Config.UsingGpg()
	if useSubprocess {
		success, err := gui.runSubprocessWithSuspense(gui.OSCommand.Cmd.NewShell(cmdObj.ToString()))
		if success && onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}
		if err := gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
			return err
		}

		return err
	} else {
		return gui.RunAndStream(cmdObj, waitingStatus, onSuccess)
	}
}

func (gui *Gui) RunAndStream(cmdObj oscommands.ICmdObj, waitingStatus string, onSuccess func() error) error {
	return gui.PopupHandler.WithWaitingStatus(waitingStatus, func() error {
		cmdObj := gui.OSCommand.Cmd.NewShell(cmdObj.ToString())
		cmdObj.AddEnvVars("TERM=dumb")
		cmdWriter := gui.getCmdWriter()
		cmd := cmdObj.GetCmd()
		cmd.Stdout = cmdWriter
		cmd.Stderr = cmdWriter

		if err := cmd.Run(); err != nil {
			if _, err := cmd.Stdout.Write([]byte(fmt.Sprintf("%s\n", style.FgRed.Sprint(err.Error())))); err != nil {
				gui.Log.Error(err)
			}
			_ = gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
			return gui.PopupHandler.Error(
				fmt.Errorf(
					gui.Tr.GitCommandFailed, gui.UserConfig.Keybinding.Universal.ExtrasMenu,
				),
			)
		}

		if onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}

		return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
	})
}
