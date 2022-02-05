package gui

import (
	"fmt"
	"os"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func (gui *Gui) submodulesRenderToMain() error {
	var task updateTask
	submodule := gui.State.Contexts.Submodules.GetSelected()
	if submodule == nil {
		task = NewRenderStringTask("No submodules")
	} else {
		prefix := fmt.Sprintf(
			"Name: %s\nPath: %s\nUrl:  %s\n\n",
			style.FgGreen.Sprint(submodule.Name),
			style.FgYellow.Sprint(submodule.Path),
			style.FgCyan.Sprint(submodule.Url),
		)

		file := gui.helpers.WorkingTree.FileForSubmodule(submodule)
		if file == nil {
			task = NewRenderStringTask(prefix)
		} else {
			cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(file, false, !file.HasUnstagedChanges && file.HasStagedChanges, gui.IgnoreWhitespaceInDiffView)
			task = NewRunCommandTaskWithPrefix(cmdObj.GetCmd(), prefix)
		}
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Submodule",
			task:  task,
		},
	})
}

func (gui *Gui) enterSubmodule(submodule *models.SubmoduleConfig) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gui.RepoPathStack.Push(wd)

	return gui.dispatchSwitchToRepo(submodule.Path, true)
}
