package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

var _ types.IController = &SwitchToSubCommitsController{}

type CanSwitchToSubCommits interface {
	types.IListContext
	GetSelectedRef() types.Ref
	ShowBranchHeadsInSubCommits() bool
}

// Not using our ListControllerTrait because our 'selected' item is not a list item
// but an attribute on it i.e. the ref of an item.
type SwitchToSubCommitsController struct {
	baseController
	*ListControllerTrait[types.Ref]
	c       *ControllerCommon
	context CanSwitchToSubCommits
}

func NewSwitchToSubCommitsController(
	c *ControllerCommon,
	context CanSwitchToSubCommits,
) *SwitchToSubCommitsController {
	return &SwitchToSubCommitsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[types.Ref](
			c,
			context,
			context.GetSelectedRef,
			func() ([]types.Ref, int, int) {
				panic("Not implemented")
			},
		),
		c:       c,
		context: context,
	}
}

func (self *SwitchToSubCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Handler:           self.viewCommits,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Description:       self.c.Tr.ViewCommits,
		},
	}

	return bindings
}

func (self *SwitchToSubCommitsController) GetOnClick() func() error {
	return self.viewCommits
}

func (self *SwitchToSubCommitsController) viewCommits() error {
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return nil
	}

	return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
		Ref:             ref,
		TitleRef:        ref.RefName(),
		Context:         self.context,
		ShowBranchHeads: self.context.ShowBranchHeadsInSubCommits(),
	})
}
