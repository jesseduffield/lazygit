package context

import (
	"os"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

// The Worktrees pane's filter matches the worktree name AND its branch, so
// searching a branch name surfaces its worktree (#5945). Pinned against the
// real FilteredList the context wires up, using the same getFilterFields
// shape NewWorktreesContext passes.
func TestWorktreesFilterMatchesBranchName(t *testing.T) {
	worktrees := []*models.Worktree{
		{Name: "frontend", Branch: "feature/login"},
		{Name: "backend", Branch: "main"},
		{Name: "docs", Branch: "docs/rewrite"},
	}

	filtered := NewFilteredList(
		func() []*models.Worktree { return worktrees },
		func(worktree *models.Worktree) []string {
			return []string{worktree.Name, worktree.Branch}
		},
	)

	getNames := func() []string {
		names := make([]string, 0)
		for _, w := range filtered.GetFilteredList() {
			names = append(names, w.Name)
		}
		return names
	}

	// By worktree name (unchanged behavior).
	filtered.SetFilter("front", false)
	assert.Equal(t, []string{"frontend"}, getNames())

	// By branch name — the #5945 ask.
	filtered.SetFilter("login", false)
	assert.Equal(t, []string{"frontend"}, getNames())

	// A branch prefix shared by name and branch fields.
	filtered.SetFilter("docs", false)
	assert.Equal(t, []string{"docs"}, getNames())

	// No match clears the list.
	filtered.SetFilter("nonexistent", false)
	assert.Empty(t, getNames())

	// Clearing the filter restores everything.
	filtered.ClearFilter()
	assert.Len(t, getFilteredListWorktrees(filtered), 3)
}

func getFilteredListWorktrees(filtered *FilteredList[*models.Worktree]) []*models.Worktree {
	return filtered.GetFilteredList()
}

// Source pin: the production context must actually wire the branch into its
// filter fields — the behavioral test above constructs its own FilteredList,
// so it cannot detect the wiring regressing.
func TestWorktreesContextWiresBranchIntoFilterFields(t *testing.T) {
	source, err := os.ReadFile("worktrees_context.go")
	if err != nil {
		t.Fatalf("read worktrees_context.go: %v", err)
	}
	assert.Contains(t, string(source),
		"return []string{worktree.Name, worktree.Branch}",
		"NewWorktreesContext must include worktree.Branch in getFilterFields (#5945)")
}
