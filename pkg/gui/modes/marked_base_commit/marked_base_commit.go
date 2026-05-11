package marked_base_commit

type MarkedBaseCommit struct {
	hash string // the hash of the commit used as a rebase base commit; empty string when unset
}

func New() MarkedBaseCommit {
	return MarkedBaseCommit{}
}

func (m *MarkedBaseCommit) Active() bool {
	return m.hash != ""
}

func (m *MarkedBaseCommit) Reset() {
	m.hash = ""
}

func (m *MarkedBaseCommit) SetHash(hash string) {
	m.hash = hash
}

func (m *MarkedBaseCommit) GetHash() string {
	return m.hash
}
