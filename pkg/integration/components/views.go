package components

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
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
	return self.newStaticViewDriver(viewName, nil, nil, nil)
}

func (self *Views) patchExplorerViewByName(viewName string) *ViewDriver {
	return self.newStaticViewDriver(
		viewName,
		func() ([]string, error) {
			ctx := self.t.gui.ContextForView(viewName).(*context.PatchExplorerContext)
			state := ctx.GetState()
			if state == nil {
				return nil, errors.New("Expected patch explorer to be activated")
			}
			selectedContent := state.PlainRenderSelected()
			// the above method returns a string with a trailing newline so we need to remove that before splitting
			selectedLines := strings.Split(strings.TrimSuffix(selectedContent, "\n"), "\n")
			return selectedLines, nil
		},
		func() (int, int, error) {
			ctx := self.t.gui.ContextForView(viewName).(*context.PatchExplorerContext)
			state := ctx.GetState()
			if state == nil {
				return 0, 0, errors.New("Expected patch explorer to be activated")
			}
			startIdx, endIdx := state.SelectedRange()
			return startIdx, endIdx, nil
		},
		func() (int, error) {
			ctx := self.t.gui.ContextForView(viewName).(*context.PatchExplorerContext)
			state := ctx.GetState()
			if state == nil {
				return 0, errors.New("Expected patch explorer to be activated")
			}
			return state.GetSelectedLineIdx(), nil
		},
	)
}

// 'static' because it'll always refer to the same view, as opposed to the 'main' view which could actually be
// one of several views, or the 'current' view which depends on focus.
func (self *Views) newStaticViewDriver(
	viewName string,
	getSelectedLinesFn func() ([]string, error),
	getSelectedLineRangeFn func() (int, int, error),
	getSelectedLineIdxFn func() (int, error),
) *ViewDriver {
	return &ViewDriver{
		context:              fmt.Sprintf("%s view", viewName),
		getView:              func() *gocui.View { return self.t.gui.View(viewName) },
		getSelectedLinesFn:   getSelectedLinesFn,
		getSelectedRangeFn:   getSelectedLineRangeFn,
		getSelectedLineIdxFn: getSelectedLineIdxFn,
		t:                    self.t,
	}
}

func (self *Views) MergeConflicts() *ViewDriver {
	viewName := "mergeConflicts"
	return self.newStaticViewDriver(
		viewName,
		func() ([]string, error) {
			ctx := self.t.gui.ContextForView(viewName).(*context.MergeConflictsContext)
			state := ctx.GetState()
			if state == nil {
				return nil, errors.New("Expected patch explorer to be activated")
			}
			selectedContent := strings.Split(state.PlainRenderSelected(), "\n")

			return selectedContent, nil
		},
		func() (int, int, error) {
			ctx := self.t.gui.ContextForView(viewName).(*context.MergeConflictsContext)
			state := ctx.GetState()
			if state == nil {
				return 0, 0, errors.New("Expected patch explorer to be activated")
			}
			startIdx, endIdx := state.GetSelectedRange()
			return startIdx, endIdx, nil
		},
		// there is no concept of a cursor in the merge conflicts panel so we just return the start of the selection
		func() (int, error) {
			ctx := self.t.gui.ContextForView(viewName).(*context.MergeConflictsContext)
			state := ctx.GetState()
			if state == nil {
				return 0, errors.New("Expected patch explorer to be activated")
			}
			startIdx, _ := state.GetSelectedRange()
			return startIdx, nil
		},
	)
}

func (self *Views) Commits() *ViewDriver {
	return self.regularView("commits")
}

func (self *Views) Files() *ViewDriver {
	return self.regularView("files")
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
