package controllers

import (
	"io/ioutil"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MergeConflictsController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &MergeConflictsController{}

func NewMergeConflictsController(
	common *controllerCommon,
) *MergeConflictsController {
	return &MergeConflictsController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *MergeConflictsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.HandleEditFile,
			Description: self.c.Tr.LcEditFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.HandleOpenFile,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.PrevBlock),
			Handler:     self.withRenderAndFocus(self.PrevConflict),
			Description: self.c.Tr.PrevConflict,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.NextBlock),
			Handler:     self.withRenderAndFocus(self.NextConflict),
			Description: self.c.Tr.NextConflict,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.PrevItem),
			Handler:     self.withRenderAndFocus(self.PrevConflictHunk),
			Description: self.c.Tr.SelectPrevHunk,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.NextItem),
			Handler:     self.withRenderAndFocus(self.NextConflictHunk),
			Description: self.c.Tr.SelectNextHunk,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevBlockAlt),
			Handler: self.withRenderAndFocus(self.PrevConflict),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextBlockAlt),
			Handler: self.withRenderAndFocus(self.NextConflict),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Handler: self.withRenderAndFocus(self.PrevConflictHunk),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextItemAlt),
			Handler: self.withRenderAndFocus(self.NextConflictHunk),
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.ScrollLeft),
			Handler:     self.withRenderAndFocus(self.HandleScrollLeft),
			Description: self.c.Tr.LcScrollLeft,
			Tag:         "navigation",
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.ScrollRight),
			Handler:     self.withRenderAndFocus(self.HandleScrollRight),
			Description: self.c.Tr.LcScrollRight,
			Tag:         "navigation",
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Undo),
			Handler:     self.withRenderAndFocus(self.HandleUndo),
			Description: self.c.Tr.LcUndo,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenMergeTool),
			Handler:     self.helpers.WorkingTree.OpenMergeTool,
			Description: self.c.Tr.LcOpenMergeTool,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.withRenderAndFocus(self.HandlePickHunk),
			Description: self.c.Tr.PickHunk,
		},
		{
			Key:         opts.GetKey(opts.Config.Main.PickBothHunks),
			Handler:     self.withRenderAndFocus(self.HandlePickAllHunks),
			Description: self.c.Tr.PickAllHunks,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.Escape,
			Description: self.c.Tr.ReturnToFilesPanel,
		},
	}

	return bindings
}

func (self *MergeConflictsController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.context().GetViewName(),
			Key:      gocui.MouseWheelUp,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.HandleScrollUp()
			},
		},
		{
			ViewName: self.context().GetViewName(),
			Key:      gocui.MouseWheelDown,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.HandleScrollDown()
			},
		},
	}
}

func (self *MergeConflictsController) HandleScrollUp() error {
	self.context().SetUserScrolling(true)
	self.context().GetViewTrait().ScrollUp(self.c.UserConfig.Gui.ScrollHeight)

	return nil
}

func (self *MergeConflictsController) HandleScrollDown() error {
	self.context().SetUserScrolling(true)
	self.context().GetViewTrait().ScrollDown(self.c.UserConfig.Gui.ScrollHeight)

	return nil
}

func (self *MergeConflictsController) Context() types.Context {
	return self.context()
}

func (self *MergeConflictsController) context() *context.MergeConflictsContext {
	return self.contexts.MergeConflicts
}

func (self *MergeConflictsController) Escape() error {
	return self.c.PushContext(self.contexts.Files)
}

func (self *MergeConflictsController) HandleEditFile() error {
	lineNumber := self.context().GetState().GetSelectedLine()
	return self.helpers.Files.EditFileAtLine(self.context().GetState().GetPath(), lineNumber)
}

func (self *MergeConflictsController) HandleOpenFile() error {
	lineNumber := self.context().GetState().GetSelectedLine()
	return self.helpers.Files.OpenFileAtLine(self.context().GetState().GetPath(), lineNumber)
}

func (self *MergeConflictsController) HandleScrollLeft() error {
	self.context().GetViewTrait().ScrollLeft()

	return nil
}

func (self *MergeConflictsController) HandleScrollRight() error {
	self.context().GetViewTrait().ScrollRight()

	return nil
}

func (self *MergeConflictsController) HandleUndo() error {
	state := self.context().GetState()

	ok := state.Undo()
	if !ok {
		return nil
	}

	self.c.LogAction("Restoring file to previous state")
	self.c.LogCommand("Undoing last conflict resolution", false)
	if err := ioutil.WriteFile(state.GetPath(), []byte(state.GetContent()), 0o644); err != nil {
		return err
	}

	return nil
}

func (self *MergeConflictsController) PrevConflictHunk() error {
	self.context().SetUserScrolling(false)
	self.context().GetState().SelectPrevConflictHunk()

	return nil
}

func (self *MergeConflictsController) NextConflictHunk() error {
	self.context().SetUserScrolling(false)
	self.context().GetState().SelectNextConflictHunk()

	return nil
}

func (self *MergeConflictsController) NextConflict() error {
	self.context().SetUserScrolling(false)
	self.context().GetState().SelectNextConflict()

	return nil
}

func (self *MergeConflictsController) PrevConflict() error {
	self.context().SetUserScrolling(false)
	self.context().GetState().SelectPrevConflict()

	return nil
}

func (self *MergeConflictsController) HandlePickHunk() error {
	return self.pickSelection(self.context().GetState().Selection())
}

func (self *MergeConflictsController) HandlePickAllHunks() error {
	return self.pickSelection(mergeconflicts.ALL)
}

func (self *MergeConflictsController) pickSelection(selection mergeconflicts.Selection) error {
	ok, err := self.resolveConflict(selection)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	if self.context().GetState().AllConflictsResolved() {
		return self.onLastConflictResolved()
	}

	return nil
}

func (self *MergeConflictsController) resolveConflict(selection mergeconflicts.Selection) (bool, error) {
	self.context().SetUserScrolling(false)

	state := self.context().GetState()

	ok, content, err := state.ContentAfterConflictResolve(selection)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	var logStr string
	switch selection {
	case mergeconflicts.TOP:
		logStr = "Picking top hunk"
	case mergeconflicts.MIDDLE:
		logStr = "Picking middle hunk"
	case mergeconflicts.BOTTOM:
		logStr = "Picking bottom hunk"
	case mergeconflicts.ALL:
		logStr = "Picking all hunks"
	}
	self.c.LogAction("Resolve merge conflict")
	self.c.LogCommand(logStr, false)
	state.PushContent(content)
	return true, ioutil.WriteFile(state.GetPath(), []byte(content), 0o644)
}

func (self *MergeConflictsController) onLastConflictResolved() error {
	// as part of refreshing files, we handle the situation where a file has had
	// its merge conflicts resolved.
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
}

func (self *MergeConflictsController) isFocused() bool {
	return self.c.CurrentContext().GetKey() == self.context().GetKey()
}

func (self *MergeConflictsController) withRenderAndFocus(f func() error) func() error {
	return self.withLock(func() error {
		if err := f(); err != nil {
			return err
		}

		return self.context().RenderAndFocus(self.isFocused())
	})
}

func (self *MergeConflictsController) withLock(f func() error) func() error {
	return func() error {
		self.context().GetMutex().Lock()
		defer self.context().GetMutex().Unlock()

		if self.context().GetState() == nil {
			return nil
		}

		return f()
	}
}
