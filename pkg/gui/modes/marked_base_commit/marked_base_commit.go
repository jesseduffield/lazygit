package marked_base_commit

type MarkedBaseCommit struct {
	sha string // the sha of the commit used as a rebase base commit; empty string when unset
}

func New() MarkedBaseCommit {
	return MarkedBaseCommit{}
}

func (m *MarkedBaseCommit) Active() bool {
	return m.sha != ""
}

func (m *MarkedBaseCommit) Reset() {
	m.sha = ""
}

func (m *MarkedBaseCommit) SetSha(sha string) {
	m.sha = sha
}

func (m *MarkedBaseCommit) GetSha() string {
	return m.sha
}
