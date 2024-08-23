package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that contain commit files.

var _ types.IController = &SwitchToDiffFilesController{}

type CanSwitchToDiffFiles interface {
	types.IListContext
	CanRebase() bool
	GetSelectedRef() types.Ref
	GetSelectedRefRangeForDiffFiles() *types.RefRange
}

// Not using our ListControllerTrait because our 'selected' item is not a list item
// but an attribute on it i.e. the ref of an item.
type SwitchToDiffFilesController struct {
	baseController
	*ListControllerTrait[types.Ref]
	c       *ControllerCommon
	context CanSwitchToDiffFiles
}

func NewSwitchToDiffFilesController(
	c *ControllerCommon,
	context CanSwitchToDiffFiles,
) *SwitchToDiffFilesController {
	return &SwitchToDiffFilesController{
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

func (self *SwitchToDiffFilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.enter,
			GetDisabledReason: self.canEnter,
			Description:       self.c.Tr.ViewItemFiles,
		},
	}

	return bindings
}

func (self *SwitchToDiffFilesController) GetOnClick() func() error {
	return func() error {
		if self.canEnter() == nil {
			return self.enter()
		}

		return nil
	}
}

func (self *SwitchToDiffFilesController) enter() error {
	ref := self.context.GetSelectedRef()
	refsRange := self.context.GetSelectedRefRangeForDiffFiles()
	commitFilesContext := self.c.Contexts().CommitFiles

	canRebase := self.context.CanRebase()
	if canRebase {
		if self.c.Modes().Diffing.Active() {
			if self.c.Modes().Diffing.Ref != ref.RefName() {
				canRebase = false
			}
		} else if refsRange != nil {
			canRebase = false
		}
	}

	commitFilesContext.ReInit(ref, refsRange)
	commitFilesContext.SetSelection(0)
	commitFilesContext.SetCanRebase(canRebase)
	commitFilesContext.SetParentContext(self.context)
	commitFilesContext.SetWindowName(self.context.GetWindowName())
	commitFilesContext.ClearSearchString()
	commitFilesContext.GetView().TitlePrefix = self.context.GetView().TitlePrefix

	if err := self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	}); err != nil {
		return err
	}

	return self.c.Context().Push(commitFilesContext)
}

func (self *SwitchToDiffFilesController) canEnter() *types.DisabledReason {
	refRange := self.context.GetSelectedRefRangeForDiffFiles()
	if refRange != nil {
		return nil
	}
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return &types.DisabledReason{Text: self.c.Tr.NoItemSelected}
	}
	if ref.RefName() == "" {
		return &types.DisabledReason{Text: self.c.Tr.SelectedItemDoesNotHaveFiles}
	}

	return nil
}
