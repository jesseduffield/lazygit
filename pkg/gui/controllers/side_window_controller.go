package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SideWindowControllerFactory struct {
	c *ControllerCommon
}

func NewSideWindowControllerFactory(common *ControllerCommon) *SideWindowControllerFactory {
	return &SideWindowControllerFactory{c: common}
}

func (self *SideWindowControllerFactory) Create(context types.Context) types.IController {
	return NewSideWindowController(self.c, context)
}

type SideWindowController struct {
	baseController
	c       *ControllerCommon
	context types.Context
}

func NewSideWindowController(
	common *ControllerCommon,
	context types.Context,
) *SideWindowController {
	return &SideWindowController{
		baseController: baseController{},
		c:              common,
		context:        context,
	}
}

func (self *SideWindowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{Key: opts.GetKey(opts.Config.Universal.PrevBlock), Modifier: gocui.ModNone, Handler: self.previousSideWindow},
		{Key: opts.GetKey(opts.Config.Universal.NextBlock), Modifier: gocui.ModNone, Handler: self.nextSideWindow},
		{Key: opts.GetKey(opts.Config.Universal.PrevBlockAlt), Modifier: gocui.ModNone, Handler: self.previousSideWindow},
		{Key: opts.GetKey(opts.Config.Universal.NextBlockAlt), Modifier: gocui.ModNone, Handler: self.nextSideWindow},
		{Key: opts.GetKey(opts.Config.Universal.PrevBlockAlt2), Modifier: gocui.ModNone, Handler: self.previousSideWindow},
		{Key: opts.GetKey(opts.Config.Universal.NextBlockAlt2), Modifier: gocui.ModNone, Handler: self.nextSideWindow},
	}
}

func (self *SideWindowController) Context() types.Context {
	return nil
}

func (self *SideWindowController) previousSideWindow() error {
	windows := self.c.Helpers().Window.SideWindows()
	currentWindow := self.c.Helpers().Window.CurrentWindow()
	var newWindow string
	if currentWindow == "" || currentWindow == windows[0] {
		newWindow = windows[len(windows)-1]
	} else {
		for i := range windows {
			if currentWindow == windows[i] {
				newWindow = windows[i-1]
				break
			}
			if i == len(windows)-1 {
				return nil
			}
		}
	}

	context := self.c.Helpers().Window.GetContextForWindow(newWindow)

	return self.c.PushContext(context)
}

func (self *SideWindowController) nextSideWindow() error {
	windows := self.c.Helpers().Window.SideWindows()
	currentWindow := self.c.Helpers().Window.CurrentWindow()
	var newWindow string
	if currentWindow == "" || currentWindow == windows[len(windows)-1] {
		newWindow = windows[0]
	} else {
		for i := range windows {
			if currentWindow == windows[i] {
				newWindow = windows[i+1]
				break
			}
			if i == len(windows)-1 {
				return nil
			}
		}
	}

	context := self.c.Helpers().Window.GetContextForWindow(newWindow)

	return self.c.PushContext(context)
}
