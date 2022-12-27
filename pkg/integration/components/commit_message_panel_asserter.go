package components

type CommitMessagePanelAsserter struct {
	input *Input
}

func (self *CommitMessagePanelAsserter) getViewAsserter() *View {
	return self.input.Views().CommitMessage()
}

// asserts on the text initially present in the prompt
func (self *CommitMessagePanelAsserter) InitialText(expected *matcher) *CommitMessagePanelAsserter {
	self.getViewAsserter().Content(expected)

	return self
}

func (self *CommitMessagePanelAsserter) Type(value string) *CommitMessagePanelAsserter {
	self.input.typeContent(value)

	return self
}

func (self *CommitMessagePanelAsserter) AddNewline() *CommitMessagePanelAsserter {
	self.input.press(self.input.keys.Universal.AppendNewline)

	return self
}

func (self *CommitMessagePanelAsserter) Clear() *CommitMessagePanelAsserter {
	panic("Clear method not yet implemented!")
}

func (self *CommitMessagePanelAsserter) Confirm() {
	self.getViewAsserter().PressEnter()
}

func (self *CommitMessagePanelAsserter) Cancel() {
	self.getViewAsserter().PressEscape()
}
