package components

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type Views struct {
	input *Input
}

// not exporting this because I want the test to always be explicit about what
// view it's dealing with.
func (self *Views) current() *View {
	return &View{
		context: "current view",
		getView: func() *gocui.View { return self.input.gui.CurrentContext().GetView() },
		input:   self.input,
	}
}

func (self *Views) Main() *View {
	return &View{
		context: "main view",
		getView: func() *gocui.View { return self.input.gui.MainView() },
		input:   self.input,
	}
}

func (self *Views) Secondary() *View {
	return &View{
		context: "secondary view",
		getView: func() *gocui.View { return self.input.gui.SecondaryView() },
		input:   self.input,
	}
}

func (self *Views) ByName(viewName string) *View {
	return &View{
		context: fmt.Sprintf("%s view", viewName),
		getView: func() *gocui.View { return self.input.gui.View(viewName) },
		input:   self.input,
	}
}

func (self *Views) Commits() *View {
	return self.ByName("commits")
}

func (self *Views) Files() *View {
	return self.ByName("files")
}

func (self *Views) Status() *View {
	return self.ByName("status")
}

func (self *Views) Submodules() *View {
	return self.ByName("submodules")
}

func (self *Views) Information() *View {
	return self.ByName("information")
}

func (self *Views) Branches() *View {
	return self.ByName("localBranches")
}

func (self *Views) RemoteBranches() *View {
	return self.ByName("remoteBranches")
}

func (self *Views) Tags() *View {
	return self.ByName("tags")
}

func (self *Views) ReflogCommits() *View {
	return self.ByName("reflogCommits")
}

func (self *Views) SubCommits() *View {
	return self.ByName("subCommits")
}

func (self *Views) CommitFiles() *View {
	return self.ByName("commitFiles")
}

func (self *Views) Stash() *View {
	return self.ByName("stash")
}

func (self *Views) Staging() *View {
	return self.ByName("staging")
}

func (self *Views) StagingSecondary() *View {
	return self.ByName("stagingSecondary")
}

func (self *Views) Menu() *View {
	return self.ByName("menu")
}

func (self *Views) Confirmation() *View {
	return self.ByName("confirmation")
}

func (self *Views) CommitMessage() *View {
	return self.ByName("commitMessage")
}

func (self *Views) Suggestions() *View {
	return self.ByName("suggestions")
}

func (self *Views) MergeConflicts() *View {
	return self.ByName("mergeConflicts")
}
