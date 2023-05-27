package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PatchExplorerControllerFactory struct {
	c *ControllerCommon
}

func NewPatchExplorerControllerFactory(c *ControllerCommon) *PatchExplorerControllerFactory {
	return &PatchExplorerControllerFactory{
		c: c,
	}
}

func (self *PatchExplorerControllerFactory) Create(context types.IPatchExplorerContext) *PatchExplorerController {
	return &PatchExplorerController{
		baseController: baseController{},
		c:              self.c,
		context:        context,
	}
}

type PatchExplorerController struct {
	baseController
	c *ControllerCommon

	context types.IPatchExplorerContext
}

func (self *PatchExplorerController) Context() types.Context {
	return self.context
}

func (self *PatchExplorerController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Handler: self.withRenderAndFocus(self.HandlePrevLine),
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.PrevItem),
			Handler: self.withRenderAndFocus(self.HandlePrevLine),
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.NextItemAlt),
			Handler: self.withRenderAndFocus(self.HandleNextLine),
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.NextItem),
			Handler: self.withRenderAndFocus(self.HandleNextLine),
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.PrevBlock),
			Handler:     self.withRenderAndFocus(self.HandlePrevHunk),
			Description: self.c.Tr.PrevHunk,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevBlockAlt),
			Handler: self.withRenderAndFocus(self.HandlePrevHunk),
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.NextBlock),
			Handler:     self.withRenderAndFocus(self.HandleNextHunk),
			Description: self.c.Tr.NextHunk,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextBlockAlt),
			Handler: self.withRenderAndFocus(self.HandleNextHunk),
		},
		{
			Key:         opts.GetKey(opts.Config.Main.ToggleDragSelect),
			Handler:     self.withRenderAndFocus(self.HandleToggleSelectRange),
			Description: self.c.Tr.ToggleDragSelect,
		},
		{
			Key:         opts.GetKey(opts.Config.Main.ToggleDragSelectAlt),
			Handler:     self.withRenderAndFocus(self.HandleToggleSelectRange),
			Description: self.c.Tr.ToggleDragSelect,
		},
		{
			Key:         opts.GetKey(opts.Config.Main.ToggleSelectHunk),
			Handler:     self.withRenderAndFocus(self.HandleToggleSelectHunk),
			Description: self.c.Tr.ToggleSelectHunk,
		},
		{
			Tag:         "navigation",
			Key:         opts.GetKey(opts.Config.Universal.PrevPage),
			Handler:     self.withRenderAndFocus(self.HandlePrevPage),
			Description: self.c.Tr.PrevPage,
		},
		{
			Tag:         "navigation",
			Key:         opts.GetKey(opts.Config.Universal.NextPage),
			Handler:     self.withRenderAndFocus(self.HandleNextPage),
			Description: self.c.Tr.NextPage,
		},
		{
			Tag:         "navigation",
			Key:         opts.GetKey(opts.Config.Universal.GotoTop),
			Handler:     self.withRenderAndFocus(self.HandleGotoTop),
			Description: self.c.Tr.GotoTop,
		},
		{
			Tag:         "navigation",
			Key:         opts.GetKey(opts.Config.Universal.GotoBottom),
			Description: self.c.Tr.GotoBottom,
			Handler:     self.withRenderAndFocus(self.HandleGotoBottom),
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.ScrollLeft),
			Handler: self.withRenderAndFocus(self.HandleScrollLeft),
		},
		{
			Tag:     "navigation",
			Key:     opts.GetKey(opts.Config.Universal.ScrollRight),
			Handler: self.withRenderAndFocus(self.HandleScrollRight),
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.withLock(self.CopySelectedToClipboard),
			Description: self.c.Tr.CopySelectedTexToClipboard,
		},
	}
}

func (self *PatchExplorerController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseLeft,
			Handler: func(opts gocui.ViewMouseBindingOpts) error {
				if self.isFocused() {
					return self.withRenderAndFocus(self.HandleMouseDown)()
				}

				return self.c.PushContext(self.context, types.OnFocusOpts{
					ClickedWindowName:  self.context.GetWindowName(),
					ClickedViewLineIdx: opts.Y,
				})
			},
		},
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModMotion,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.withRenderAndFocus(self.HandleMouseDrag)()
			},
		},
	}
}

func (self *PatchExplorerController) HandlePrevLine() error {
	self.context.GetState().CycleSelection(false)

	return nil
}

func (self *PatchExplorerController) HandleNextLine() error {
	self.context.GetState().CycleSelection(true)

	return nil
}

func (self *PatchExplorerController) HandlePrevHunk() error {
	self.context.GetState().CycleHunk(false)

	return nil
}

func (self *PatchExplorerController) HandleNextHunk() error {
	self.context.GetState().CycleHunk(true)

	return nil
}

func (self *PatchExplorerController) HandleToggleSelectRange() error {
	self.context.GetState().ToggleSelectRange()

	return nil
}

func (self *PatchExplorerController) HandleToggleSelectHunk() error {
	self.context.GetState().ToggleSelectHunk()

	return nil
}

func (self *PatchExplorerController) HandleScrollLeft() error {
	self.context.GetViewTrait().ScrollLeft()

	return nil
}

func (self *PatchExplorerController) HandleScrollRight() error {
	self.context.GetViewTrait().ScrollRight()

	return nil
}

func (self *PatchExplorerController) HandlePrevPage() error {
	self.context.GetState().SetLineSelectMode()
	self.context.GetState().AdjustSelectedLineIdx(-self.context.GetViewTrait().PageDelta())

	return nil
}

func (self *PatchExplorerController) HandleNextPage() error {
	self.context.GetState().SetLineSelectMode()
	self.context.GetState().AdjustSelectedLineIdx(self.context.GetViewTrait().PageDelta())

	return nil
}

func (self *PatchExplorerController) HandleGotoTop() error {
	self.context.GetState().SelectTop()

	return nil
}

func (self *PatchExplorerController) HandleGotoBottom() error {
	self.context.GetState().SelectBottom()

	return nil
}

func (self *PatchExplorerController) HandleMouseDown() error {
	self.context.GetState().SelectNewLineForRange(self.context.GetViewTrait().SelectedLineIdx())

	return nil
}

func (self *PatchExplorerController) HandleMouseDrag() error {
	self.context.GetState().SelectLine(self.context.GetViewTrait().SelectedLineIdx())

	return nil
}

func (self *PatchExplorerController) CopySelectedToClipboard() error {
	selected := self.context.GetState().PlainRenderSelected()

	self.c.LogAction(self.c.Tr.Actions.CopySelectedTextToClipboard)
	if err := self.c.OS().CopyToClipboard(selected); err != nil {
		return self.c.Error(err)
	}

	return nil
}

func (self *PatchExplorerController) isFocused() bool {
	return self.c.CurrentContext().GetKey() == self.context.GetKey()
}

func (self *PatchExplorerController) withRenderAndFocus(f func() error) func() error {
	return self.withLock(func() error {
		if err := f(); err != nil {
			return err
		}

		return self.context.RenderAndFocus(self.isFocused())
	})
}

func (self *PatchExplorerController) withLock(f func() error) func() error {
	return func() error {
		self.context.GetMutex().Lock()
		defer self.context.GetMutex().Unlock()

		if self.context.GetState() == nil {
			return nil
		}

		return f()
	}
}
