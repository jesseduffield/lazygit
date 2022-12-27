package components

type CommitMessagePanelAsserter struct {
	t *TestDriver
}

func (self *CommitMessagePanelAsserter) getViewAsserter() *View {
	return self.t.Views().CommitMessage()
}

// asserts on the text initially present in the prompt
func (self *CommitMessagePanelAsserter) InitialText(expected *matcher) *CommitMessagePanelAsserter {
	self.getViewAsserter().Content(expected)

	return self
}

func (self *CommitMessagePanelAsserter) Type(value string) *CommitMessagePanelAsserter {
	self.t.typeContent(value)

	return self
}

func (self *CommitMessagePanelAsserter) AddNewline() *CommitMessagePanelAsserter {
	self.t.press(self.t.keys.Universal.AppendNewline)

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
