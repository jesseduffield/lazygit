package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SpiceStacksController struct {
	baseController
	*ListControllerTrait[*models.SpiceStackItem]
	c            *ControllerCommon
	hasRefreshed bool // Track if we've refreshed the data
}

var _ types.IController = &SpiceStacksController{}

func NewSpiceStacksController(
	c *ControllerCommon,
) *SpiceStacksController {
	return &SpiceStacksController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().SpiceStacks,
			c.Contexts().SpiceStacks.GetSelected,
			c.Contexts().SpiceStacks.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SpiceStacksController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Handler: self.HandlePrevLine,
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.PrevItem),
			Handler: self.HandlePrevLine,
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.NextItemAlt),
			Handler: self.HandleNextLine,
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.NextItem),
			Handler: self.HandleNextLine,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.checkout),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.newBranch,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "New branch",
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.delete),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Delete branch",
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranch),
			Handler:           self.withItem(self.restack),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Restack",
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey("R"),
			Handler:         self.restackAll,
			Description:     "Restack all",
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CreatePullRequest),
			Handler:           self.withItem(self.submit),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Submit PR",
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey("O"),
			Handler:         self.submitAll,
			Description:     "Submit all",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("<c-u>"),
			Handler:         self.navigateUp,
			Description:     "Up stack",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("<c-d>"),
			Handler:         self.navigateDown,
			Description:     "Down stack",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("<c-U>"),
			Handler:         self.navigateTop,
			Description:     "Top of stack",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("<c-D>"),
			Handler:         self.navigateBottom,
			Description:     "Bottom of stack",
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey("<c-j>"),
			Handler:           self.withItem(self.moveBranchDown),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Move branch down in stack",
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey("<c-k>"),
			Handler:           self.withItem(self.moveBranchUp),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Move branch up in stack",
			DisplayOnScreen:   true,
		},
		{
			Key:         opts.GetKey("l"),
			Handler:     self.toggleLogFormat,
			Description: self.c.Tr.ToggleSpiceLogFormat,
			Tooltip:     self.c.Tr.ToggleSpiceLogFormatTooltip,
		},
	}

	return bindings
}

func (self *SpiceStacksController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		// Ensure we start on a branch, not a commit
		self.ensureValidSelection()
	}
}

// HandleNextLine moves to the next branch, skipping commits
func (self *SpiceStacksController) HandleNextLine() error {
	return self.handleLineChange(1)
}

// HandlePrevLine moves to the previous branch, skipping commits
func (self *SpiceStacksController) HandlePrevLine() error {
	return self.handleLineChange(-1)
}

// handleLineChange navigates to the next/prev branch, skipping commits
func (self *SpiceStacksController) handleLineChange(delta int) error {
	items := self.c.Model().SpiceStackItems
	if len(items) == 0 {
		return nil
	}

	currentIdx := self.context().GetSelectedLineIdx()
	newIdx := currentIdx + delta

	// Find the next non-commit item in the given direction
	for newIdx >= 0 && newIdx < len(items) {
		if !items[newIdx].IsCommit {
			self.context().SetSelection(newIdx)
			self.context().HandleFocus(types.OnFocusOpts{})
			return nil
		}
		newIdx += delta
	}

	// No valid item found - stay at current position
	return nil
}

// ensureValidSelection makes sure we're not selecting a commit
func (self *SpiceStacksController) ensureValidSelection() {
	items := self.c.Model().SpiceStackItems
	if len(items) == 0 {
		return
	}

	currentIdx := self.context().GetSelectedLineIdx()
	if currentIdx < 0 || currentIdx >= len(items) {
		// Find first branch
		for i, item := range items {
			if !item.IsCommit {
				self.context().SetSelection(i)
				return
			}
		}
		return
	}

	// If current selection is a commit, move to nearest branch
	if items[currentIdx].IsCommit {
		// Try moving forward first
		for i := currentIdx + 1; i < len(items); i++ {
			if !items[i].IsCommit {
				self.context().SetSelection(i)
				return
			}
		}
		// Try moving backward
		for i := currentIdx - 1; i >= 0; i-- {
			if !items[i].IsCommit {
				self.context().SetSelection(i)
				return
			}
		}
	}
}

func (self *SpiceStacksController) checkout(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return nil // Commits are not interactive
	}
	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	if err := self.c.Git().Branch.Checkout(item.Name, git_commands.CheckoutOptions{Force: false}); err != nil {
		return err
	}
	self.hasRefreshed = false // Reset so data refreshes on next focus
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) restack(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return nil // Commits are not interactive
	}
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack(item.Name)
		if err != nil {
			return err
		}
		self.hasRefreshed = false
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) restackAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack("")
		if err != nil {
			return err
		}
		self.hasRefreshed = false
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) submit(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return nil // Commits are not interactive
	}
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit(item.Name)
		if err != nil {
			return err
		}
		self.hasRefreshed = false
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) submitAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit("")
		if err != nil {
			return err
		}
		self.hasRefreshed = false
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) navigateUp() error {
	if err := self.c.Git().Spice.NavigateUp(); err != nil {
		return err
	}
	self.hasRefreshed = false
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) navigateDown() error {
	if err := self.c.Git().Spice.NavigateDown(); err != nil {
		return err
	}
	self.hasRefreshed = false
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) navigateTop() error {
	if err := self.c.Git().Spice.NavigateTop(); err != nil {
		return err
	}
	self.hasRefreshed = false
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) navigateBottom() error {
	if err := self.c.Git().Spice.NavigateBottom(); err != nil {
		return err
	}
	self.hasRefreshed = false
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) newBranch() error {
	self.c.Prompt(types.PromptOpts{
		Title: "Branch name:",
		HandleConfirm: func(branchName string) error {
			if err := self.c.Git().Spice.CreateBranch(branchName); err != nil {
				return err
			}
			self.hasRefreshed = false
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) delete(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return nil // Commits are not interactive
	}
	self.c.Confirm(types.ConfirmOpts{
		Title:  "Delete branch",
		Prompt: "Are you sure you want to delete this branch from the stack?",
		HandleConfirm: func() error {
			if err := self.c.Git().Spice.DeleteBranch(item.Name); err != nil {
				return err
			}
			self.hasRefreshed = false
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.SPICE_STACKS}})
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) moveBranchUp(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return nil // Commits are not interactive
	}
	if err := self.c.Git().Spice.MoveBranchUp(item.Name); err != nil {
		return err
	}
	self.hasRefreshed = false
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) moveBranchDown(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return nil // Commits are not interactive
	}
	if err := self.c.Git().Spice.MoveBranchDown(item.Name); err != nil {
		return err
	}
	self.hasRefreshed = false
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) toggleLogFormat() error {
	currentFormat := self.c.UserConfig().Git.Spice.LogFormat

	if currentFormat == "long" {
		self.c.UserConfig().Git.Spice.LogFormat = "short"
	} else {
		self.c.UserConfig().Git.Spice.LogFormat = "long"
	}

	// Refresh the spice stacks view
	self.c.Refresh(types.RefreshOptions{
		Mode:  types.ASYNC,
		Scope: []types.RefreshableView{types.SPICE_STACKS},
	})
	return nil
}

func (self *SpiceStacksController) context() *context.SpiceStacksContext {
	return self.c.Contexts().SpiceStacks
}

func (self *SpiceStacksController) withItem(f func(item *models.SpiceStackItem) error) func() error {
	return func() error {
		item := self.context().GetSelected()
		if item == nil {
			return nil
		}
		return f(item)
	}
}

func (self *SpiceStacksController) singleItemSelected() func() *types.DisabledReason {
	return func() *types.DisabledReason {
		if self.context().GetSelected() == nil {
			return &types.DisabledReason{Text: "No item selected"}
		}
		return nil
	}
}

func (self *SpiceStacksController) require(conditions ...func() *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		for _, condition := range conditions {
			if reason := condition(); reason != nil {
				return reason
			}
		}
		return nil
	}
}
