package diffing

// if ref is blank we're not diffing anything
type Diffing struct {
	Ref     string
	Reverse bool
}

func New() Diffing {
	return Diffing{}
}

func (m *Diffing) Active() bool {
	return m.Ref != ""
}
