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
