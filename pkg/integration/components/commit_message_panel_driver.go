package components

type CommitMessagePanelDriver struct {
	t *TestDriver
}

func (self *CommitMessagePanelDriver) getViewDriver() *ViewDriver {
	return self.t.Views().CommitMessage()
}

// asserts on the text initially present in the prompt
func (self *CommitMessagePanelDriver) InitialText(expected *matcher) *CommitMessagePanelDriver {
	self.getViewDriver().Content(expected)

	return self
}

func (self *CommitMessagePanelDriver) Type(value string) *CommitMessagePanelDriver {
	self.t.typeContent(value)

	return self
}

func (self *CommitMessagePanelDriver) AddNewline() *CommitMessagePanelDriver {
	self.t.press(self.t.keys.Universal.AppendNewline)

	return self
}

func (self *CommitMessagePanelDriver) Clear() *CommitMessagePanelDriver {
	panic("Clear method not yet implemented!")
}

func (self *CommitMessagePanelDriver) Confirm() {
	self.getViewDriver().PressEnter()
}

func (self *CommitMessagePanelDriver) Cancel() {
	self.getViewDriver().PressEscape()
}
