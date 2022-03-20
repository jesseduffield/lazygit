package git_commands

import (
	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/generics/slices"
	"github.com/sirupsen/logrus"
)

// although the typical terms in a git bisect are 'bad' and 'good', they're more
// generally known as 'new' and 'old'. Semi-recently git allowed the user to define
// their own terms e.g. when you want to used 'fixed', 'unfixed' in the event
// that you're looking for a commit that fixed a bug.

// Git bisect only keeps track of a single 'bad' commit. Once you pick a commit
// that's older than the current bad one, it forgets about the previous one. On
// the other hand, it does keep track of all the good and skipped commits.

type BisectInfo struct {
	log *logrus.Entry

	// tells us whether all our git bisect files are there meaning we're in bisect mode.
	// Doesn't necessarily mean that we've actually picked a good/bad commit yet.
	started bool

	// this is the ref you started the commit from
	start string // this will always be defined

	// these will be defined if we've started
	newTerm string // 'bad' by default
	oldTerm string // 'good' by default

	// map of commit sha's to their status
	statusMap map[string]BisectStatus

	// the sha of the commit that's under test
	current string
}

type BisectStatus int

const (
	BisectStatusOld BisectStatus = iota
	BisectStatusNew
	BisectStatusSkipped
)

// null object pattern
func NewNullBisectInfo() *BisectInfo {
	return &BisectInfo{started: false}
}

func (self *BisectInfo) GetNewSha() string {
	for sha, status := range self.statusMap {
		if status == BisectStatusNew {
			return sha
		}
	}

	return ""
}

func (self *BisectInfo) GetCurrentSha() string {
	return self.current
}

func (self *BisectInfo) GetStartSha() string {
	return self.start
}

func (self *BisectInfo) Status(commitSha string) (BisectStatus, bool) {
	status, ok := self.statusMap[commitSha]
	return status, ok
}

func (self *BisectInfo) NewTerm() string {
	return self.newTerm
}

func (self *BisectInfo) OldTerm() string {
	return self.oldTerm
}

// this is for when we have called `git bisect start`. It does not
// mean that we have actually started narrowing things down or selecting good/bad commits
func (self *BisectInfo) Started() bool {
	return self.started
}

// this is where we have both a good and bad revision and we're actually
// starting to narrow things down
func (self *BisectInfo) Bisecting() bool {
	if !self.Started() {
		return false
	}

	if self.GetNewSha() == "" {
		return false
	}

	return slices.Contains(maps.Values(self.statusMap), BisectStatusOld)
}
