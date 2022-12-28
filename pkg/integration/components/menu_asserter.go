package components

type MenuAsserter struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *MenuAsserter) getViewAsserter() *View {
	return self.t.Views().Menu()
}

// asserts that the popup has the expected title
func (self *MenuAsserter) Title(expected *matcher) *MenuAsserter {
	self.getViewAsserter().Title(expected)

	self.hasCheckedTitle = true

	return self
}

func (self *MenuAsserter) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.getViewAsserter().PressEnter()
}

func (self *MenuAsserter) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewAsserter().PressEscape()
}

func (self *MenuAsserter) Select(option *matcher) *MenuAsserter {
	self.getViewAsserter().NavigateToListItem(option)

	return self
}

func (self *MenuAsserter) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.t.Fail("You must check the title of a menu popup by calling Title() before calling Confirm()/Cancel().")
	}
}
