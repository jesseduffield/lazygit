package gui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

type InputImpl struct {
	g    *gocui.Gui
	keys config.KeybindingConfig
}

var _ types.Input = &InputImpl{}

func (self *InputImpl) PushKeys(keyStrs ...string) {
	for _, keyStr := range keyStrs {
		self.pushKey(keyStr)
	}
}

func (self *InputImpl) pushKey(keyStr string) {
	key := keybindings.GetKey(keyStr)

	var r rune
	var tcellKey tcell.Key
	switch v := key.(type) {
	case rune:
		r = v
		tcellKey = tcell.KeyRune
	case gocui.Key:
		tcellKey = tcell.Key(v)
	}

	self.g.ReplayedEvents.Keys <- gocui.NewTcellKeyEventWrapper(
		tcell.NewEventKey(tcellKey, r, tcell.ModNone),
		0,
	)
}

func (self *InputImpl) SwitchToStatusWindow() {
	self.pushKey(self.keys.Universal.JumpToBlock[0])
}

func (self *InputImpl) SwitchToFilesWindow() {
	self.pushKey(self.keys.Universal.JumpToBlock[1])
}

func (self *InputImpl) SwitchToBranchesWindow() {
	self.pushKey(self.keys.Universal.JumpToBlock[2])
}

func (self *InputImpl) SwitchToCommitsWindow() {
	self.pushKey(self.keys.Universal.JumpToBlock[3])
}

func (self *InputImpl) SwitchToStashWindow() {
	self.pushKey(self.keys.Universal.JumpToBlock[4])
}

func (self *InputImpl) Type(content string) {
	for _, char := range content {
		self.pushKey(string(char))
	}
}

func (self *InputImpl) Confirm() {
	self.pushKey(self.keys.Universal.Confirm)
}

func (self *InputImpl) Cancel() {
	self.pushKey(self.keys.Universal.Return)
}

func (self *InputImpl) Select() {
	self.pushKey(self.keys.Universal.Select)
}

func (self *InputImpl) NextItem() {
	self.pushKey(self.keys.Universal.NextItem)
}

func (self *InputImpl) PreviousItem() {
	self.pushKey(self.keys.Universal.PrevItem)
}

func (self *InputImpl) Wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
