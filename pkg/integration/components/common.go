package components

// for running common actions
type Common struct {
	t *TestDriver
}

func (self *Common) ContinueMerge() {
	self.t.GlobalPress(self.t.keys.Universal.CreateRebaseOptionsMenu)

	self.t.ExpectPopup().Menu().
		Title(Equals("Rebase options")).
		Select(Contains("continue")).
		Confirm()
}

func (self *Common) ContinueRebase() {
	self.ContinueMerge()
}

func (self *Common) AcknowledgeConflicts() {
	self.t.ExpectPopup().Menu().
		Title(Equals("Conflicts!")).
		Select(Contains("View conflicts")).
		Confirm()
}

func (self *Common) ContinueOnConflictsResolved() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("Continue")).
		Content(Contains("All merge conflicts resolved. Continue?")).
		Confirm()
}

func (self *Common) ConfirmDiscardLines() {
	self.t.ExpectPopup().Confirmation().
		Title(Equals("Discard change")).
		Content(Contains("Are you sure you want to discard this change")).
		Confirm()
}

func (self *Common) SelectPatchOption(matcher *TextMatcher) {
	self.t.GlobalPress(self.t.keys.Universal.CreatePatchOptionsMenu)

	self.t.ExpectPopup().Menu().Title(Equals("Patch options")).Select(matcher).Confirm()
}
