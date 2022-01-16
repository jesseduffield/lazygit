package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// list panel functions

func (gui *Gui) getSelectedRemote() *models.Remote {
	selectedLine := gui.State.Panels.Remotes.SelectedLineIdx
	if selectedLine == -1 || len(gui.State.Remotes) == 0 {
		return nil
	}

	return gui.State.Remotes[selectedLine]
}

func (gui *Gui) remotesRenderToMain() error {
	var task updateTask
	remote := gui.getSelectedRemote()
	if remote == nil {
		task = NewRenderStringTask("No remotes")
	} else {
		task = NewRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", style.FgGreen.Sprint(remote.Name), strings.Join(remote.Urls, "\n")))
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Remote",
			task:  task,
		},
	})
}

func (gui *Gui) refreshRemotes() error {
	prevSelectedRemote := gui.getSelectedRemote()

	remotes, err := gui.git.Loaders.Remotes.GetRemotes()
	if err != nil {
		return gui.c.Error(err)
	}

	gui.State.Remotes = remotes

	// we need to ensure our selected remote branches aren't now outdated
	if prevSelectedRemote != nil && gui.State.RemoteBranches != nil {
		// find remote now
		for _, remote := range remotes {
			if remote.Name == prevSelectedRemote.Name {
				gui.State.RemoteBranches = remote.Branches
			}
		}
	}

	return gui.c.PostRefreshUpdate(gui.mustContextForContextKey(types.ContextKey(gui.Views.Branches.Context)))
}
