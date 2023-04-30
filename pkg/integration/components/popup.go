package components

type Popup struct {
	t *TestDriver
}

func (self *Popup) Confirmation() *ConfirmationDriver {
	self.inConfirm()

	return &ConfirmationDriver{t: self.t}
}

func (self *Popup) inConfirm() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected confirmation popup to be focused"
	})
}

func (self *Popup) Prompt() *PromptDriver {
	self.inPrompt()

	return &PromptDriver{t: self.t}
}

func (self *Popup) inPrompt() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && currentView.Editable, "Expected prompt popup to be focused"
	})
}

func (self *Popup) Alert() *AlertDriver {
	self.inAlert()

	return &AlertDriver{t: self.t}
}

func (self *Popup) inAlert() {
	// basically the same thing as a confirmation popup with the current implementation
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "confirmation" && !currentView.Editable, "Expected alert popup to be focused"
	})
}

func (self *Popup) Menu() *MenuDriver {
	self.inMenu()

	return &MenuDriver{t: self.t}
}

func (self *Popup) inMenu() {
	self.t.assertWithRetries(func() (bool, string) {
		return self.t.gui.CurrentContext().GetView().Name() == "menu", "Expected popup menu to be focused"
	})
}

func (self *Popup) CommitMessagePanel() *CommitMessagePanelDriver {
	self.inCommitMessagePanel()

	return &CommitMessagePanelDriver{t: self.t}
}

func (self *Popup) CommitDescriptionPanel() *CommitMessagePanelDriver {
	self.inCommitDescriptionPanel()

	return &CommitMessagePanelDriver{t: self.t}
}

func (self *Popup) inCommitMessagePanel() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "commitMessage", "Expected commit message panel to be focused"
	})
}

func (self *Popup) inCommitDescriptionPanel() {
	self.t.assertWithRetries(func() (bool, string) {
		currentView := self.t.gui.CurrentContext().GetView()
		return currentView.Name() == "commitDescription", "Expected commit description panel to be focused"
	})
}
