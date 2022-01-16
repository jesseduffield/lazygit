package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (self *Gui) getSelectedTag() *models.Tag {
	selectedLine := self.State.Panels.Tags.SelectedLineIdx
	if selectedLine == -1 || len(self.State.Tags) == 0 {
		return nil
	}

	return self.State.Tags[selectedLine]
}

func (self *Gui) tagsRenderToMain() error {
	var task updateTask
	tag := self.getSelectedTag()
	if tag == nil {
		task = NewRenderStringTask("No tags")
	} else {
		cmdObj := self.git.Branch.GetGraphCmdObj(tag.Name)
		task = NewRunCommandTask(cmdObj.GetCmd())
	}

	return self.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Tag",
			task:  task,
		},
	})
}

// this is a controller: it can't access tags directly. Or can it? It should be able to get but not set. But that's exactly what I'm doing here, setting it. but through a mutator which encapsulates the event.
func (self *Gui) refreshTags() error {
	tags, err := self.git.Loaders.Tags.GetTags()
	if err != nil {
		return self.c.Error(err)
	}

	self.State.Tags = tags

	return self.postRefreshUpdate(self.State.Contexts.Tags)
}
