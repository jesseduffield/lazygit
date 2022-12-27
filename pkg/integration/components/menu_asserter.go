package components

type MenuAsserter struct {
	assert          *Assert
	input           *Input
	hasCheckedTitle bool
}

func (self *MenuAsserter) getViewAsserter() *Views {
	return self.assert.Views().ByName("menu")
}

// asserts that the popup has the expected title
func (self *MenuAsserter) Title(expected *matcher) *MenuAsserter {
	self.getViewAsserter().Title(expected)

	self.hasCheckedTitle = true

	return self
}

func (self *MenuAsserter) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.input.Confirm()
}

func (self *MenuAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.input.Press(self.input.keys.Universal.Return)
}

func (self *MenuAsserter) Select(option *matcher) *MenuAsserter {
	self.input.NavigateToListItem(option)

	return self
}

func (self *MenuAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.assert.Fail("You must check the title of a menu popup by calling Title() before calling Confirm()/Cancel().")
	}
}
