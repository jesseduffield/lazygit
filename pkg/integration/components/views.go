package components

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type Views struct {
	t *TestDriver
}

// not exporting this because I want the test to always be explicit about what
// view it's dealing with.
func (self *Views) current() *ViewDriver {
	return &ViewDriver{
		context: "current view",
		getView: func() *gocui.View { return self.t.gui.CurrentContext().GetView() },
		t:       self.t,
	}
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

func (self *Views) byName(viewName string) *ViewDriver {
	return &ViewDriver{
		context: fmt.Sprintf("%s view", viewName),
		getView: func() *gocui.View { return self.t.gui.View(viewName) },
		t:       self.t,
	}
}

func (self *Views) Commits() *ViewDriver {
	return self.byName("commits")
}

func (self *Views) Files() *ViewDriver {
	return self.byName("files")
}

func (self *Views) Status() *ViewDriver {
	return self.byName("status")
}

func (self *Views) Submodules() *ViewDriver {
	return self.byName("submodules")
}

func (self *Views) Information() *ViewDriver {
	return self.byName("information")
}

func (self *Views) AppStatus() *ViewDriver {
	return self.byName("appStatus")
}

func (self *Views) Branches() *ViewDriver {
	return self.byName("localBranches")
}

func (self *Views) RemoteBranches() *ViewDriver {
	return self.byName("remoteBranches")
}

func (self *Views) Tags() *ViewDriver {
	return self.byName("tags")
}

func (self *Views) ReflogCommits() *ViewDriver {
	return self.byName("reflogCommits")
}

func (self *Views) SubCommits() *ViewDriver {
	return self.byName("subCommits")
}

func (self *Views) CommitFiles() *ViewDriver {
	return self.byName("commitFiles")
}

func (self *Views) Stash() *ViewDriver {
	return self.byName("stash")
}

func (self *Views) Staging() *ViewDriver {
	return self.byName("staging")
}

func (self *Views) StagingSecondary() *ViewDriver {
	return self.byName("stagingSecondary")
}

func (self *Views) Menu() *ViewDriver {
	return self.byName("menu")
}

func (self *Views) Confirmation() *ViewDriver {
	return self.byName("confirmation")
}

func (self *Views) CommitMessage() *ViewDriver {
	return self.byName("commitMessage")
}

func (self *Views) Suggestions() *ViewDriver {
	return self.byName("suggestions")
}

func (self *Views) MergeConflicts() *ViewDriver {
	return self.byName("mergeConflicts")
}
