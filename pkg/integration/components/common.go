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

func (self *Common) AbortRebase() {
	self.t.GlobalPress(self.t.keys.Universal.CreateRebaseOptionsMenu)

	self.t.ExpectPopup().Menu().
		Title(Equals("Rebase options")).
		Select(Contains("abort")).
		Confirm()
}

func (self *Common) AbortMerge() {
	self.t.GlobalPress(self.t.keys.Universal.CreateRebaseOptionsMenu)

	self.t.ExpectPopup().Menu().
		Title(Equals("Merge options")).
		Select(Contains("abort")).
		Confirm()
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

func (self *Common) ResetBisect() {
	self.t.Views().Commits().
		Focus().
		Press(self.t.keys.Commits.ViewBisectOptions).
		Tap(func() {
			self.t.ExpectPopup().Menu().
				Title(Equals("Bisect")).
				Select(Contains("Reset bisect")).
				Confirm()

			self.t.ExpectPopup().Confirmation().
				Title(Equals("Reset 'git bisect'")).
				Content(Contains("Are you sure you want to reset 'git bisect'?")).
				Confirm()
		})
}

func (self *Common) ResetCustomPatch() {
	self.t.GlobalPress(self.t.keys.Universal.CreatePatchOptionsMenu)

	self.t.ExpectPopup().Menu().
		Title(Equals("Patch options")).
		Select(Contains("Reset patch")).
		Confirm()
}
