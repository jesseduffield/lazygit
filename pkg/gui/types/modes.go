package types

import (
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/filtering"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/marked_base_commit"
)

type Modes struct {
	Filtering filtering.Filtering

	// Shared between all worktrees of the same repo (see gui.SharedRepoState).
	// Mutate it through this pointer, but never replace it, otherwise it is no
	// longer shared.
	CherryPicking *cherrypicking.CherryPicking

	Diffing          diffing.Diffing
	MarkedBaseCommit marked_base_commit.MarkedBaseCommit
}
