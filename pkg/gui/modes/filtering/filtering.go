package filtering

type Filtering struct {
	path string // the filename that gets passed to git log
}

func New(path string) Filtering {
	return Filtering{path: path}
}

func (m *Filtering) Active() bool {
	return m.path != ""
}

func (m *Filtering) Reset() {
	m.path = ""
}

func (m *Filtering) SetPath(path string) {
	m.path = path
}

func (m *Filtering) GetPath() string {
	return m.path
}
