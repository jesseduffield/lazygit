package components

type CommitMessagePanelAsserter struct {
	assert *Assert
	input  *Input
}

func (self *CommitMessagePanelAsserter) getViewAsserter() *ViewAsserter {
	return self.assert.View("commitMessage")
}

// asserts on the text initially present in the prompt
func (self *CommitMessagePanelAsserter) InitialText(expected *matcher) *CommitMessagePanelAsserter {
	self.getViewAsserter().Content(expected)

	return self
}

func (self *CommitMessagePanelAsserter) Type(value string) *CommitMessagePanelAsserter {
	self.input.Type(value)

	return self
}

func (self *CommitMessagePanelAsserter) AddNewline() *CommitMessagePanelAsserter {
	self.input.Press(self.input.keys.Universal.AppendNewline)

	return self
}

func (self *CommitMessagePanelAsserter) Clear() *CommitMessagePanelAsserter {
	panic("Clear method not yet implemented!")
}

func (self *CommitMessagePanelAsserter) Confirm() {
	self.input.Confirm()
}

func (self *CommitMessagePanelAsserter) Cancel() {
	self.input.Press(self.input.keys.Universal.Return)
}
