package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
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
	// Focusing by keyboard doesn't point at any particular line: in a diff view
	// we start at the first change block (like entering the staging view), in a
	// non-diff view we show no selection and the user just scrolls. Clicking does
	// point at a line, so it selects that line instead (focusMainView's
	// clickedLineIdx).
	return self.focusMainView(self.c.Contexts().Normal, -1)
}

func (self *SwitchToFocusedMainViewController) focusMainView(mainViewContext *context.MainContext, clickedLineIdx int) error {
	mainViewContext.ClearSearchString()
	self.c.Context().Push(mainViewContext, types.OnFocusOpts{})

	if !sidePanelShowsDiff(self.context) {
		// Non-diff main content (e.g. a branch's commit log): focus only, no
		// selection, since there's nothing to act on.
		return nil
	}

	// A click points at a specific line (clickedLineIdx); keyboard focus (-1) starts at
	// the first change block. Either way, if the pager produced an unresolvable diff
	// this re-renders it raw first (see establishFocusedDiffSelection).
	establishFocusedDiffSelection(self.c, mainViewContext, clickedLineIdx)

	// The inclusion gutter is refreshed by the main view's focus handler (it's shown
	// only while focused, so it tracks focus changes there).
	return nil
}
