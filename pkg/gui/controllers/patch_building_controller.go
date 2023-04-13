package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type PatchBuildingController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &PatchBuildingController{}

func NewPatchBuildingController(
	common *controllerCommon,
) *PatchBuildingController {
	return &PatchBuildingController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *PatchBuildingController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.OpenFile,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.EditFile,
			Description: self.c.Tr.LcEditFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.ToggleSelectionAndRefresh,
			Description: self.c.Tr.ToggleSelectionForPatch,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.Escape,
			Description: self.c.Tr.ExitCustomPatchBuilder,
		},
	}
}

func (self *PatchBuildingController) Context() types.Context {
	return self.contexts.CustomPatchBuilder
}

func (self *PatchBuildingController) context() types.IPatchExplorerContext {
	return self.contexts.CustomPatchBuilder
}

func (self *PatchBuildingController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{}
}

func (self *PatchBuildingController) OpenFile() error {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	path := self.contexts.CommitFiles.GetSelectedPath()

	if path == "" {
		return nil
	}

	return self.helpers.Files.OpenFile(path)
}

func (self *PatchBuildingController) EditFile() error {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	path := self.contexts.CommitFiles.GetSelectedPath()

	if path == "" {
		return nil
	}

	lineNumber := self.context().GetState().CurrentLineNumber()
	return self.helpers.Files.EditFileAtLine(path, lineNumber)
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

	toggleFunc := self.git.Patch.PatchBuilder.AddFileLineRange
	filename := self.contexts.CommitFiles.GetSelectedPath()
	if filename == "" {
		return nil
	}

	state := self.context().GetState()

	includedLineIndices, err := self.git.Patch.PatchBuilder.GetFileIncLineIndices(filename)
	if err != nil {
		return err
	}
	currentLineIsStaged := lo.Contains(includedLineIndices, state.GetSelectedLineIdx())
	if currentLineIsStaged {
		toggleFunc = self.git.Patch.PatchBuilder.RemoveFileLineRange
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
	return self.helpers.PatchBuilding.Escape()
}
