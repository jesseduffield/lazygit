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
func (self *Views) current() *View {
	return &View{
		context: "current view",
		getView: func() *gocui.View { return self.t.gui.CurrentContext().GetView() },
		t:       self.t,
	}
}

func (self *Views) Main() *View {
	return &View{
		context: "main view",
		getView: func() *gocui.View { return self.t.gui.MainView() },
		t:       self.t,
	}
}

func (self *Views) Secondary() *View {
	return &View{
		context: "secondary view",
		getView: func() *gocui.View { return self.t.gui.SecondaryView() },
		t:       self.t,
	}
}

func (self *Views) byName(viewName string) *View {
	return &View{
		context: fmt.Sprintf("%s view", viewName),
		getView: func() *gocui.View { return self.t.gui.View(viewName) },
		t:       self.t,
	}
}

func (self *Views) Commits() *View {
	return self.byName("commits")
}

func (self *Views) Files() *View {
	return self.byName("files")
}

func (self *Views) Status() *View {
	return self.byName("status")
}

func (self *Views) Submodules() *View {
	return self.byName("submodules")
}

func (self *Views) Information() *View {
	return self.byName("information")
}

func (self *Views) Branches() *View {
	return self.byName("localBranches")
}

func (self *Views) RemoteBranches() *View {
	return self.byName("remoteBranches")
}

func (self *Views) Tags() *View {
	return self.byName("tags")
}

func (self *Views) ReflogCommits() *View {
	return self.byName("reflogCommits")
}

func (self *Views) SubCommits() *View {
	return self.byName("subCommits")
}

func (self *Views) CommitFiles() *View {
	return self.byName("commitFiles")
}

func (self *Views) Stash() *View {
	return self.byName("stash")
}

func (self *Views) Staging() *View {
	return self.byName("staging")
}

func (self *Views) StagingSecondary() *View {
	return self.byName("stagingSecondary")
}

func (self *Views) Menu() *View {
	return self.byName("menu")
}

func (self *Views) Confirmation() *View {
	return self.byName("confirmation")
}

func (self *Views) CommitMessage() *View {
	return self.byName("commitMessage")
}

func (self *Views) Suggestions() *View {
	return self.byName("suggestions")
}

func (self *Views) MergeConflicts() *View {
	return self.byName("mergeConflicts")
}
