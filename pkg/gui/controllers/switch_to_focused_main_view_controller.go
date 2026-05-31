package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that can focus their main view.

var _ types.IController = &SwitchToFocusedMainViewController{}

type SwitchToFocusedMainViewController struct {
	baseController
	c       *ControllerCommon
	context types.Context
}

func NewSwitchToFocusedMainViewController(
	c *ControllerCommon,
	context types.Context,
) *SwitchToFocusedMainViewController {
	return &SwitchToFocusedMainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
	}
}

func (self *SwitchToFocusedMainViewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:        opts.GetKeys(opts.Config.Universal.FocusMainView),
			Handler:     self.handleFocusMainView,
			Description: self.c.Tr.FocusMainView,
			Tag:         "global",
		},
	}

	return bindings
}

func (self *SwitchToFocusedMainViewController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    "main",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FocusedView: self.context.GetViewName(),
		},
		{
			ViewName:    "secondary",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickSecondary,
			FocusedView: self.context.GetViewName(),
		},
	}
}

func (self *SwitchToFocusedMainViewController) Context() types.Context {
	return self.context
}

func (self *SwitchToFocusedMainViewController) onClickMain(opts gocui.ViewMouseBindingOpts) error {
	return self.focusMainView(self.c.Contexts().Normal, opts.Y)
}

func (self *SwitchToFocusedMainViewController) onClickSecondary(opts gocui.ViewMouseBindingOpts) error {
	return self.focusMainView(self.c.Contexts().NormalSecondary, opts.Y)
}

func (self *SwitchToFocusedMainViewController) handleFocusMainView() error {
	// Focusing by keyboard doesn't point at any particular line, so we don't
	// show a selection; the user is free to scroll. Clicking does point at a
	// line, so it selects it (see focusMainView's clickedLineIdx).
	return self.focusMainView(self.c.Contexts().Normal, -1)
}

func (self *SwitchToFocusedMainViewController) focusMainView(mainViewContext types.Context, clickedLineIdx int) error {
	if context, ok := mainViewContext.(types.ISearchableContext); ok {
		context.ClearSearchString()
	}
	self.c.Context().Push(mainViewContext, types.OnFocusOpts{})
	if clickedLineIdx >= 0 {
		showSelectionAtLine(mainViewContext.GetView(), clickedLineIdx)
	}
	return nil
}
