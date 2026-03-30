package models

// SpiceStackItem is the display model for a branch in the stack tree
type SpiceStackItem struct {
	Name          string
	Current       bool
	Depth         int // indentation level in tree
	SiblingIndex  int // position among siblings (0 = topmost)
	NeedsRestack  bool
	PRNumber      string // e.g. "#123"
	PRStatus      string // "open", "closed", "merged"
	Ahead         int
	Behind        int
	IsCommit      bool   // marks commit items vs branch items
	CommitSha     string // for commit items
	CommitSubject string // for commit items
}

func (s *SpiceStackItem) ID() string {
	return s.Name
}

func (s *SpiceStackItem) URN() string {
	return "spice-stack-" + s.ID()
}

func (s *SpiceStackItem) RefName() string {
	return s.Name
}

func (s *SpiceStackItem) ShortRefName() string {
	return s.RefName()
}

func (s *SpiceStackItem) FullRefName() string {
	return "refs/heads/" + s.Name
}

func (s *SpiceStackItem) ParentRefName() string {
	return s.RefName() + "^"
}

func (s *SpiceStackItem) Description() string {
	return s.Name
}

// SpiceBranchJSON matches gs log long --json output
type SpiceBranchJSON struct {
	Name    string            `json:"name"`
	Current bool              `json:"current,omitempty"`
	Down    *SpiceDownJSON    `json:"down,omitempty"`
	Ups     []SpiceUpJSON     `json:"ups,omitempty"`
	Change  *SpiceChange      `json:"change,omitempty"`
	Push    *SpicePush        `json:"push,omitempty"`
	Commits []SpiceCommitJSON `json:"commits,omitempty"`
}

type SpiceCommitJSON struct {
	Sha     string `json:"sha"`
	Subject string `json:"subject"`
}

type SpiceDownJSON struct {
	Name         string `json:"name"`
	NeedsRestack bool   `json:"needsRestack,omitempty"`
}

type SpiceUpJSON struct {
	Name string `json:"name"`
}

type SpiceChange struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status,omitempty"`
}

type SpicePush struct {
	Ahead     int  `json:"ahead"`
	Behind    int  `json:"behind"`
	NeedsPush bool `json:"needsPush,omitempty"`
}
