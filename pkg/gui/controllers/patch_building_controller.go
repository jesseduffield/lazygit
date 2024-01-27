package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type PatchBuildingController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &PatchBuildingController{}

func NewPatchBuildingController(
	c *ControllerCommon,
) *PatchBuildingController {
	return &PatchBuildingController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *PatchBuildingController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.OpenFile,
			Description: self.c.Tr.OpenFile,
			Tooltip:     self.c.Tr.OpenFileTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.EditFile,
			Description: self.c.Tr.EditFile,
			Tooltip:     self.c.Tr.EditFileTooltip,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.Select),
			Handler:         self.ToggleSelectionAndRefresh,
			Description:     self.c.Tr.ToggleSelectionForPatch,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.Escape,
			Description: self.c.Tr.ExitCustomPatchBuilder,
		},
	}
}

func (self *PatchBuildingController) Context() types.Context {
	return self.c.Contexts().CustomPatchBuilder
}

func (self *PatchBuildingController) context() types.IPatchExplorerContext {
	return self.c.Contexts().CustomPatchBuilder
}

func (self *PatchBuildingController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{}
}

func (self *PatchBuildingController) GetOnFocus() func(types.OnFocusOpts) error {
	return func(opts types.OnFocusOpts) error {
		// no need to change wrap on the secondary view because it can't be interacted with
		self.c.Views().PatchBuilding.Wrap = false

		return self.c.Helpers().PatchBuilding.RefreshPatchBuildingPanel(opts)
	}
}

func (self *PatchBuildingController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(opts types.OnFocusLostOpts) error {
		self.c.Views().PatchBuilding.Wrap = true

		if self.c.Git().Patch.PatchBuilder.IsEmpty() {
			self.c.Git().Patch.PatchBuilder.Reset()
		}

		return nil
	}
}

func (self *PatchBuildingController) OpenFile() error {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	path := self.c.Contexts().CommitFiles.GetSelectedPath()

	if path == "" {
		return nil
	}

	return self.c.Helpers().Files.OpenFile(path)
}

func (self *PatchBuildingController) EditFile() error {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	path := self.c.Contexts().CommitFiles.GetSelectedPath()

	if path == "" {
		return nil
	}

	lineNumber := self.context().GetState().CurrentLineNumber()
	return self.c.Helpers().Files.EditFileAtLine(path, lineNumber)
}

func (self *PatchBuildingController) ToggleSelectionAndRefresh() error {
	if err := self.toggleSelection(); err != nil {
		return err
	}

	return self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.PATCH_BUILDING, types.COMMIT_FILES},
	})
}

func (self *PatchBuildingController) toggleSelection() error {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	toggleFunc := self.c.Git().Patch.PatchBuilder.AddFileLineRange
	filename := self.c.Contexts().CommitFiles.GetSelectedPath()
	if filename == "" {
		return nil
	}

	state := self.context().GetState()

	includedLineIndices, err := self.c.Git().Patch.PatchBuilder.GetFileIncLineIndices(filename)
	if err != nil {
		return err
	}
	currentLineIsStaged := lo.Contains(includedLineIndices, state.GetSelectedLineIdx())
	if currentLineIsStaged {
		toggleFunc = self.c.Git().Patch.PatchBuilder.RemoveFileLineRange
	}

	// add range of lines to those set for the file
	firstLineIdx, lastLineIdx := state.SelectedRange()

	if err := toggleFunc(filename, firstLineIdx, lastLineIdx); err != nil {
		// might actually want to return an error here
		self.c.Log.Error(err)
	}

	if state.SelectingRange() {
		state.SetLineSelectMode()
	}

	return nil
}

func (self *PatchBuildingController) Escape() error {
	context := self.c.Contexts().CustomPatchBuilder
	state := context.GetState()

	if state.SelectingRange() || state.SelectingHunk() {
		state.SetLineSelectMode()
		return self.c.PostRefreshUpdate(context)
	}

	return self.c.Helpers().PatchBuilding.Escape()
}
