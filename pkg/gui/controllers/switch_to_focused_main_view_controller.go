package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
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
			Key:         opts.GetKey(opts.Config.Universal.FocusMainView),
			Handler:     self.handleFocusMainView,
			Description: self.c.Tr.FocusMainView,
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
	return self.focusMainView("main", opts.Y)
}

func (self *SwitchToFocusedMainViewController) onClickSecondary(opts gocui.ViewMouseBindingOpts) error {
	return self.focusMainView("secondary", opts.Y)
}

func (self *SwitchToFocusedMainViewController) handleFocusMainView() error {
	return self.focusMainView("main", -1)
}

func (self *SwitchToFocusedMainViewController) focusMainView(mainViewName string, clickedViewLineIdx int) error {
	mainViewContext := self.c.Helpers().Window.GetContextForWindow(mainViewName)
	mainViewContext.SetParentContext(self.context)
	if context, ok := mainViewContext.(types.ISearchableContext); ok {
		context.ClearSearchString()
	}
	onFocusOpts := types.OnFocusOpts{ClickedWindowName: mainViewName}
	if clickedViewLineIdx >= 0 {
		onFocusOpts.ClickedViewLineIdx = clickedViewLineIdx
	} else {
		mainView := mainViewContext.GetView()
		lineIdx := mainView.OriginY() + mainView.Height()/2
		onFocusOpts.ClickedViewLineIdx = lo.Clamp(lineIdx, 0, mainView.LinesHeight()-1)
	}
	self.c.Context().Push(mainViewContext, onFocusOpts)
	return nil
}
