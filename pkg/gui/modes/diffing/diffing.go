package diffing

// if ref is blank we're not diffing anything
type Diffing struct {
	Ref     string
	Reverse bool
}

func New() Diffing {
	return Diffing{}
}

func (self *Diffing) Active() bool {
	return self.Ref != ""
}

// GetFromAndReverseArgsForDiff tells us the from and reverse args to be used in a diff command.
// If we're not in diff mode we'll end up with the equivalent of a `git show` i.e `git diff blah^..blah`.
func (self *Diffing) GetFromAndReverseArgsForDiff(from string) (string, bool) {
	reverse := false

	if self.Active() {
		reverse = self.Reverse
		from = self.Ref
	}

	return from, reverse
}
