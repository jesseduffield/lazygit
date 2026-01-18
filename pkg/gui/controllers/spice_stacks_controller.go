package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SpiceStacksController struct {
	baseController
	*ListControllerTrait[*models.SpiceStackItem]
	c           *ControllerCommon
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
	}

	return bindings
}

func (self *SpiceStacksController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		// Only refresh once when first focusing the tab
		if !self.hasRefreshed {
			self.hasRefreshed = true
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		}
	}
}

func (self *SpiceStacksController) checkout(item *models.SpiceStackItem) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	return self.c.Helpers().Refs.CheckoutRef(item.Name, types.CheckoutRefOptions{})
}

func (self *SpiceStacksController) restack(item *models.SpiceStackItem) error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack(item.Name)
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	})
}

func (self *SpiceStacksController) restackAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack("")
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	})
}

func (self *SpiceStacksController) submit(item *models.SpiceStackItem) error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit(item.Name)
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	})
}

func (self *SpiceStacksController) submitAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit("")
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	})
}

func (self *SpiceStacksController) navigateUp() error {
	if err := self.c.Git().Spice.NavigateUp(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *SpiceStacksController) navigateDown() error {
	if err := self.c.Git().Spice.NavigateDown(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *SpiceStacksController) navigateTop() error {
	if err := self.c.Git().Spice.NavigateTop(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *SpiceStacksController) navigateBottom() error {
	if err := self.c.Git().Spice.NavigateBottom(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *SpiceStacksController) newBranch() error {
	self.c.Prompt(types.PromptOpts{
		Title: "Branch name:",
		HandleConfirm: func(branchName string) error {
			if err := self.c.Git().Spice.CreateBranch(branchName); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) delete(item *models.SpiceStackItem) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  "Delete branch",
		Prompt: "Are you sure you want to delete this branch from the stack?",
		HandleConfirm: func() error {
			if err := self.c.Git().Spice.DeleteBranch(item.Name); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) moveBranchUp(item *models.SpiceStackItem) error {
	if err := self.c.Git().Spice.MoveBranchUp(item.Name); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *SpiceStacksController) moveBranchDown(item *models.SpiceStackItem) error {
	if err := self.c.Git().Spice.MoveBranchDown(item.Name); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
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
