package gui

import "github.com/jesseduffield/lazygit/pkg/commands/models"

// this file is to put things where it's not obvious where they belong while this refactor takes place

func (gui *Gui) getSuggestedRemote() string {
	remotes := gui.State.Remotes

	return getSuggestedRemote(remotes)
}

func getSuggestedRemote(remotes []*models.Remote) string {
	if len(remotes) == 0 {
		return "origin"
	}

	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote.Name
		}
	}

	return remotes[0].Name
}
