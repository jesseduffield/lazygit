package components

// for running common actions
type Common struct {
	t *TestDriver
}

func (self *Common) ContinueMerge() {
	self.t.GlobalPress(self.t.keys.Universal.CreateRebaseOptionsMenu)

	self.t.ExpectPopup().Menu().
		Title(Equals("Rebase Options")).
		Select(Contains("continue")).
		Confirm()
}

func (self *Common) ContinueRebase() {
	self.ContinueMerge()
}

func (self *Common) AcknowledgeConflicts() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("Auto-merge failed")).
		Content(Contains("Conflicts!")).
		Confirm()
}

func (self *Common) ContinueOnConflictsResolved() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("continue")).
		Content(Contains("all merge conflicts resolved. Continue?")).
		Confirm()
}

func (self *Common) ConfirmDiscardLines() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("Unstage lines")).
		Content(Contains("Are you sure you want to delete the selected lines")).
		Confirm()
}

func (self *Common) SelectPatchOption(matcher *Matcher) {
	self.t.GlobalPress(self.t.keys.Universal.CreatePatchOptionsMenu)

	self.t.ExpectPopup().Menu().Title(Equals("Patch Options")).Select(matcher).Confirm()
}
