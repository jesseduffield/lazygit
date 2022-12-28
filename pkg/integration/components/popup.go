package components

type Popup struct {
	t *TestDriver
}

func (self *Popup) Confirmation() *ConfirmationAsserter {
	self.inConfirm()

	return &ConfirmationAsserter{t: self.t}
}

func (self *Popup) inConfirm() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected confirmation popup to be focused"
	})
}

func (self *Popup) Prompt() *PromptAsserter {
	self.inPrompt()

	return &PromptAsserter{t: self.t}
}

func (self *Popup) inPrompt() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && currentView.Editable, "Expected prompt popup to be focused"
	})
}

func (self *Popup) Alert() *AlertAsserter {
	self.inAlert()

	return &AlertAsserter{t: self.t}
}

func (self *Popup) inAlert() {
	// basically the same thing as a confirmation popup with the current implementation
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected alert popup to be focused"
	})
}

func (self *Popup) Menu() *MenuAsserter {
	self.inMenu()

	return &MenuAsserter{t: self.t}
}

func (self *Popup) inMenu() {
	self.t.assertWithRetries(func() (bool, string) {
		return self.t.gui.CurrentContext().GetView().Name() == "menu", "Expected popup menu to be focused"
	})
}

func (self *Popup) CommitMessagePanel() *CommitMessagePanelAsserter {
	self.inCommitMessagePanel()

	return &CommitMessagePanelAsserter{t: self.t}
}

func (self *Popup) inCommitMessagePanel() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "commitMessage", "Expected commit message panel to be focused"
	})
}
