package components

// for running common actions
type Actions struct {
	t *TestDriver
}

func (self *Actions) ContinueMerge() {
	self.t.Views().current().Press(self.t.keys.Universal.CreateRebaseOptionsMenu)

	self.t.ExpectPopup().Menu().
		Title(Equals("Rebase Options")).
		Select(Contains("continue")).
		Confirm()
}

func (self *Actions) ContinueRebase() {
	self.ContinueMerge()
}

func (self *Actions) AcknowledgeConflicts() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("Auto-merge failed")).
		Content(Contains("Conflicts!")).
		Confirm()
}

func (self *Actions) ContinueOnConflictsResolved() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("continue")).
		Content(Contains("all merge conflicts resolved. Continue?")).
		Confirm()
}

func (self *Actions) ConfirmDiscardLines() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("Unstage lines")).
		Content(Contains("Are you sure you want to delete the selected lines")).
		Confirm()
}
