package components

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type Views struct {
	t *TestDriver
}

func (self *Views) Main() *ViewDriver {
	return &ViewDriver{
		context: "main view",
		getView: func() *gocui.View { return self.t.gui.MainView() },
		t:       self.t,
	}
}

func (self *Views) Secondary() *ViewDriver {
	return &ViewDriver{
		context: "secondary view",
		getView: func() *gocui.View { return self.t.gui.SecondaryView() },
		t:       self.t,
	}
}

func (self *Views) regularView(viewName string) *ViewDriver {
	return &ViewDriver{
		context: fmt.Sprintf("%s view", viewName),
		getView: func() *gocui.View { return self.t.gui.View(viewName) },
		t:       self.t,
	}
}

func (self *Views) patchExplorerViewByName(viewName string) *ViewDriver {
	return self.regularView(viewName)
}

func (self *Views) MergeConflicts() *ViewDriver {
	return self.regularView("mergeConflicts")
}

func (self *Views) Commits() *ViewDriver {
	return self.regularView("commits")
}

func (self *Views) Files() *ViewDriver {
	return self.regularView("files")
}

func (self *Views) Worktrees() *ViewDriver {
	return self.regularView("worktrees")
}

func (self *Views) Status() *ViewDriver {
	return self.regularView("status")
}

func (self *Views) Submodules() *ViewDriver {
	return self.regularView("submodules")
}

func (self *Views) Information() *ViewDriver {
	return self.regularView("information")
}

func (self *Views) AppStatus() *ViewDriver {
	return self.regularView("appStatus")
}

func (self *Views) Branches() *ViewDriver {
	return self.regularView("localBranches")
}

func (self *Views) Remotes() *ViewDriver {
	return self.regularView("remotes")
}

func (self *Views) RemoteBranches() *ViewDriver {
	return self.regularView("remoteBranches")
}

func (self *Views) Tags() *ViewDriver {
	return self.regularView("tags")
}

func (self *Views) ReflogCommits() *ViewDriver {
	return self.regularView("reflogCommits")
}

func (self *Views) SubCommits() *ViewDriver {
	return self.regularView("subCommits")
}

func (self *Views) CommitFiles() *ViewDriver {
	return self.regularView("commitFiles")
}

func (self *Views) Stash() *ViewDriver {
	return self.regularView("stash")
}

func (self *Views) Staging() *ViewDriver {
	return self.patchExplorerViewByName("staging")
}

func (self *Views) StagingSecondary() *ViewDriver {
	return self.patchExplorerViewByName("stagingSecondary")
}

func (self *Views) PatchBuilding() *ViewDriver {
	return self.patchExplorerViewByName("patchBuilding")
}

func (self *Views) PatchBuildingSecondary() *ViewDriver {
	// this is not a patch explorer view because you can't actually focus it: it
	// just renders content
	return self.regularView("patchBuildingSecondary")
}

func (self *Views) Menu() *ViewDriver {
	return self.regularView("menu")
}

func (self *Views) Confirmation() *ViewDriver {
	return self.regularView("confirmation")
}

func (self *Views) CommitMessage() *ViewDriver {
	return self.regularView("commitMessage")
}

func (self *Views) CommitDescription() *ViewDriver {
	return self.regularView("commitDescription")
}

func (self *Views) Suggestions() *ViewDriver {
	return self.regularView("suggestions")
}

func (self *Views) Search() *ViewDriver {
	return self.regularView("search")
}

func (self *Views) Tooltip() *ViewDriver {
	return self.regularView("tooltip")
}

func (self *Views) Options() *ViewDriver {
	return self.regularView("options")
}
