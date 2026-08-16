package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
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
	// Usually the main pane, but the content can be in the secondary one alone: a file
	// with nothing but staged changes shows them there.
	mainViewContext := self.c.Contexts().Normal
	if self.c.State().GetRepoState().GetMainPanes() == types.SecondaryPaneOnly {
		mainViewContext = self.c.Contexts().NormalSecondary
	}
	return self.focusMainView(mainViewContext, -1)
}

func (self *SwitchToFocusedMainViewController) focusMainView(mainViewContext *context.MainContext, clickedLineIdx int) error {
	mainViewContext.ClearSearchString()
	self.c.Context().Push(mainViewContext, types.OnFocusOpts{})

	if _, ok := self.context.(types.DiffMainViewContext); !ok {
		return nil
	}

	// The diff on screen was produced for reading, and the renderer that produced it may
	// have laid it out in a way that says nothing about which line of which file each row
	// is. Now that the user wants to act on it, it is re-rendered as git's own diff — the
	// panel below decides that for itself, from the same question — and the selection
	// goes on that instead of on rows we can't place.
	if self.c.Helpers().DiffLine.MainViewDiffMode() == git_commands.DiffModeRaw {
		self.c.Helpers().DiffLine.RenderFocusedMainViewAgain(mainViewContext.GetView(), self.context, func() {
			self.c.Helpers().DiffLine.EstablishSelection(mainViewContext, clickedLineIdx)
		})
		return nil
	}

	self.c.Helpers().DiffLine.EstablishSelection(mainViewContext, clickedLineIdx)
	return nil
}
