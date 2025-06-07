package controllers

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
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

// Not using our ListControllerTrait because we have our own way of working with
// range selections that's different from ListControllerTrait's
type SwitchToDiffFilesController struct {
	baseController
	c       *ControllerCommon
	context CanSwitchToDiffFiles
}

func NewSwitchToDiffFilesController(
	c *ControllerCommon,
	context CanSwitchToDiffFiles,
) *SwitchToDiffFilesController {
	return &SwitchToDiffFilesController{
		baseController: baseController{},
		c:              c,
		context:        context,
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

func (self *SwitchToDiffFilesController) GetOnClickFocusedMainView() func(mainViewName string, clickedLineIdx int) error {
	return func(mainViewName string, clickedLineIdx int) error {
		clickedFile, line, ok := self.c.Helpers().Staging.GetFileAndLineForClickedDiffLine(mainViewName, clickedLineIdx)
		if !ok {
			return nil
		}

		if err := self.enter(); err != nil {
			return err
		}

		context := self.c.Contexts().CommitFiles
		var node *filetree.CommitFileNode

		relativePath, err := filepath.Rel(self.c.Git().RepoPaths.RepoPath(), clickedFile)
		if err != nil {
			return err
		}
		relativePath = "./" + relativePath
		context.CommitFileTreeViewModel.ExpandToPath(relativePath)
		self.c.PostRefreshUpdate(context)

		idx, ok := context.CommitFileTreeViewModel.GetIndexForPath(relativePath)
		if !ok {
			return nil
		}

		context.SetSelectedLineIdx(idx)
		context.GetViewTrait().FocusPoint(
			context.ModelIndexToViewIndex(idx))
		node = context.GetSelected()
		self.c.Contexts().CustomPatchBuilder.SetParentContext(self.context)
		return self.c.Helpers().CommitFiles.EnterCommitFile(node, types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: line, ClickedViewRealLineIdx: line})
	}
}

func (self *SwitchToDiffFilesController) Context() types.Context {
	return self.context
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

	self.c.Context().Push(commitFilesContext, types.OnFocusOpts{})
	return nil
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
