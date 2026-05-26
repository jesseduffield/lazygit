package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type MainViewController struct {
	baseController
	c *ControllerCommon

	context      *context.MainContext
	otherContext *context.MainContext
}

var _ types.IController = &MainViewController{}

func NewMainViewController(
	c *ControllerCommon,
	context *context.MainContext,
	otherContext *context.MainContext,
) *MainViewController {
	return &MainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
		otherContext:   otherContext,
	}
}

func (self *MainViewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	// When a selection is shown, we surface the bindings that act on it
	// (enter to dive into staging, escape to hide the selection).
	selectionShown := self.context.GetView().Highlight

	var enterDescription string
	if selectionShown {
		enterDescription = self.c.Tr.EnterStaging
	}

	return []*types.Binding{
		{
			Keys:            opts.GetKeys(opts.Config.Universal.TogglePanel),
			Handler:         self.togglePanel,
			Description:     self.c.Tr.ToggleStagingView,
			Tooltip:         self.c.Tr.ToggleStagingViewTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Return),
			Handler:         self.escape,
			Description:     self.c.Tr.ExitFocusedMainView,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Select),
			Handler:         self.toggleSelection,
			Description:     self.c.Tr.ToggleSelectionInFocusedMainView,
			DisplayOnScreen: !selectionShown,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.GoInto),
			Handler:         self.enter,
			Description:     enterDescription,
			DisplayOnScreen: selectionShown,
		},
		{
			// overriding this because we want to read all of the task's output before we start searching
			Keys:        opts.GetKeys(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.StartSearch,
			Tag:         "navigation",
		},
	}
}

func (self *MainViewController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInAlreadyFocusedView,
			FocusedView: self.context.GetViewName(),
		},
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInOtherViewOfMainViewPair,
			FocusedView: self.otherContext.GetViewName(),
		},
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
}

// Transient focus shifts (popups, search) leave HighlightInactive=true on our
// view (set by ContextMgr.Activate when a different view becomes current). Our
// context's highlightOnFocus is false, so SimpleContext.HandleFocus never
// resets it. Reset it here on the way back in, so that if we still hold a
// selection it's drawn as active. The flag is a no-op when Highlight is false.
func (self *MainViewController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		self.context.GetView().HighlightInactive = false
	}
}

func (self *MainViewController) togglePanel() error {
	if self.otherContext.GetView().Visible {
		self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	}

	return nil
}

func (self *MainViewController) escape() error {
	v := self.context.GetView()
	if v.Highlight {
		v.Highlight = false
		return nil
	}
	self.c.Context().Pop()
	return nil
}

func (self *MainViewController) toggleSelection() error {
	v := self.context.GetView()
	if v.Highlight {
		v.Highlight = false
		return nil
	}
	v.Highlight = true
	v.HighlightInactive = false
	lineIdx := v.OriginY() + v.InnerHeight()/2
	lineIdx = lo.Clamp(lineIdx, 0, v.ViewLinesHeight()-1)
	v.FocusPoint(0, lineIdx, false)
	return nil
}

func (self *MainViewController) enter() error {
	if !self.context.GetView().Highlight {
		return nil
	}
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(
			self.context.GetViewName(), self.context.GetView().SelectedLineIdx())
	}
	return nil
}

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	if self.context.GetView().Highlight && !opts.IsDoubleClick {
		return nil
	}

	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(self.context.GetViewName(), opts.Y)
	}
	return nil
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Push(self.context, types.OnFocusOpts{
		ClickedWindowName:  self.context.GetWindowName(),
		ClickedViewLineIdx: opts.Y,
	})

	return nil
}

func (self *MainViewController) openSearch() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				return self.c.Helpers().Search.OpenSearchPrompt(self.context)
			})
		})
	}

	return nil
}
