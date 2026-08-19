package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// CommitDiffActions implements what a panel showing a commit's diff offers on that diff
// in the focused main view. Five panels do: the commit files panel shows the diff of one
// file of a commit, and the commits, sub-commits, stash and reflog panels the whole diff
// of whatever they have selected. They all offer the same thing and differ only in which
// diff they show, so they share this, each saying which diff that is.
type CommitDiffActions struct {
	c *ControllerCommon

	// The panel this belongs to, and what it is showing the diff of — nil when it has
	// nothing selected, and so no diff.
	panel  types.Context
	target func() *commitDiffTarget
}

// commitDiffTarget is the diff a panel is showing: the two ends of it, and whether it
// belongs to a commit lazygit may rewrite.
type commitDiffTarget struct {
	from      string
	to        string
	canRebase bool
}

var _ types.FocusedMainViewDiffSource = &CommitDiffActions{}

func NewCommitDiffActions(
	c *ControllerCommon, panel types.Context, target func() *commitDiffTarget,
) *CommitDiffActions {
	return &CommitDiffActions{c: c, panel: panel, target: target}
}

// PlainDiff hands out the diff the panel is showing, for the given files — the same
// diff as in the main view, only without the commit's message and stat above it.
func (self *CommitDiffActions) PlainDiff(_ types.DiffPaneContext, paths []string) string {
	target := self.target()
	if target == nil {
		return ""
	}
	return self.c.Helpers().Diff.PlainDiffBetweenRefs(target.from, target.to, paths)
}
