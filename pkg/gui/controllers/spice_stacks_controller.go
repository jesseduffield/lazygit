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
	c *ControllerCommon
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
			Handler:           self.withItem(self.restack),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Restack",
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey("<c-r>"),
			Handler:         self.restackAll,
			Description:     "Restack all",
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey("S"),
			Handler:           self.withItem(self.submit),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       "Submit PR",
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey("<c-s>"),
			Handler:         self.submitAll,
			Description:     "Submit all",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("u"),
			Handler:         self.navigateUp,
			Description:     "Up stack",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("d"),
			Handler:         self.navigateDown,
			Description:     "Down stack",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("U"),
			Handler:         self.navigateTop,
			Description:     "Top of stack",
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("D"),
			Handler:         self.navigateBottom,
			Description:     "Bottom of stack",
			DisplayOnScreen: true,
		},
	}

	return bindings
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
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *SpiceStacksController) restackAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack("")
		if err != nil {
			return err
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *SpiceStacksController) submit(item *models.SpiceStackItem) error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit(item.Name)
		if err != nil {
			return err
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *SpiceStacksController) submitAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit("")
		if err != nil {
			return err
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *SpiceStacksController) navigateUp() error {
	if err := self.c.Git().Spice.NavigateUp(); err != nil {
		return err
	}
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *SpiceStacksController) navigateDown() error {
	if err := self.c.Git().Spice.NavigateDown(); err != nil {
		return err
	}
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *SpiceStacksController) navigateTop() error {
	if err := self.c.Git().Spice.NavigateTop(); err != nil {
		return err
	}
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *SpiceStacksController) navigateBottom() error {
	if err := self.c.Git().Spice.NavigateBottom(); err != nil {
		return err
	}
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
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
