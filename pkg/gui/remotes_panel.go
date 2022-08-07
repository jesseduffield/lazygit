package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// list panel functions

func (gui *Gui) remotesRenderToMain() error {
	var task types.UpdateTask
	remote := gui.State.Contexts.Remotes.GetSelected()
	if remote == nil {
		task = types.NewRenderStringTask("No remotes")
	} else {
		task = types.NewRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", style.FgGreen.Sprint(remote.Name), strings.Join(remote.Urls, "\n")))
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: "Remote",
			Task:  task,
		},
	})
}
