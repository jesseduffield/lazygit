package components

import (
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

// through this struct we assert on the state of the lazygit gui

type Assert struct {
	input *Input
	gui   integrationTypes.GuiDriver
	*assertionHelper
}

func NewAssert(gui integrationTypes.GuiDriver) *Assert {
	return &Assert{gui: gui}
}

// for making assertions on lazygit views
func (self *Assert) Views() *Views {
	return &Views{assert: self, input: self.input}
}

// for making assertions on the lazygit model
func (self *Assert) Model() *Model {
	return &Model{assertionHelper: self.assertionHelper, gui: self.gui}
}

// for making assertions on the file system
func (self *Assert) FileSystem() *FileSystem {
	return &FileSystem{assertionHelper: self.assertionHelper}
}

// for when you just want to fail the test yourself.
// This runs callbacks to ensure we render the error after closing the gui.
func (self *Assert) Fail(message string) {
	self.assertionHelper.fail(message)
}
