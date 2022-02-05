package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

// list panel functions

func (gui *Gui) remotesRenderToMain() error {
	var task updateTask
	remote := gui.State.Contexts.Remotes.GetSelected()
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
