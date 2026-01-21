package presentation

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildTreePrefix(t *testing.T) {
	// Test data matching the current repo structure from gs log short output
	// Expected output:
	//         ┏━□ spice-stacks-alt-long-nav
	//         ┃ ┏━□ spice-stack-refactor
	//         ┣━┻□ spice-stacks-nav-bug
	//       ┏━┻■ spice-stacks-rendering ◀
	//     ┏━┻□ spice-stacks-sort-order
	//   ┏━┻□ spice-stacks-long-desc
	// ┏━┻□ claude/add-spice-stacks-tab-xljBG
	// master

	// Items are in reversed order (leaves first)
	items := []*models.SpiceStackItem{
		{Name: "spice-stacks-alt-long-nav", Depth: 5, SiblingIndex: 0},
		{Name: "spice-stack-refactor", Depth: 6, SiblingIndex: 0},
		{Name: "spice-stacks-nav-bug", Depth: 5, SiblingIndex: 1},
		{Name: "spice-stacks-rendering", Depth: 4, SiblingIndex: 0, Current: true},
		{Name: "spice-stacks-sort-order", Depth: 3, SiblingIndex: 0},
		{Name: "spice-stacks-long-desc", Depth: 2, SiblingIndex: 0},
		{Name: "claude/add-spice-stacks-tab-xljBG", Depth: 1, SiblingIndex: 0},
		{Name: "master", Depth: 0, SiblingIndex: 0},
	}

	tests := []struct {
		idx      int
		expected string
	}{
		{0, "        ┌──◯ "},   // spice-stacks-alt-long-nav
		{1, "        │ ┌──◯ "}, // spice-stack-refactor
		{2, "        ├─┴◯ "},   // spice-stacks-nav-bug
		{3, "      ┌─┴● "},     // spice-stacks-rendering (current)
		{4, "    ┌─┴◯ "},       // spice-stacks-sort-order
		{5, "  ┌─┴◯ "},         // spice-stacks-long-desc
		{6, "┌─┴◯ "},           // claude/add-spice-stacks-tab-xljBG
		{7, ""},                // master: no prefix
	}

	for _, tt := range tests {
		result := buildTreePrefix(items[tt.idx], tt.idx, items)
		// Remove any ANSI color codes if present
		result = stripAnsi(result)
		assert.Equal(t, tt.expected, result, "Item %d (%s) tree prefix mismatch", tt.idx, items[tt.idx].Name)
	}
}

func TestBuildTreePrefixSimpleLinear(t *testing.T) {
	// Test simple linear stack: A -> B -> C
	items := []*models.SpiceStackItem{
		{Name: "C", Depth: 3, SiblingIndex: 0},
		{Name: "B", Depth: 2, SiblingIndex: 0},
		{Name: "A", Depth: 1, SiblingIndex: 0},
		{Name: "main", Depth: 0, SiblingIndex: 0},
	}

	tests := []struct {
		idx      int
		expected string
	}{
		{0, "    ┌──◯ "}, // C
		{1, "  ┌─┴◯ "},   // B
		{2, "┌─┴◯ "},     // A
		{3, ""},          // main: no prefix
	}

	for _, tt := range tests {
		result := buildTreePrefix(items[tt.idx], tt.idx, items)
		result = stripAnsi(result)
		assert.Equal(t, tt.expected, result, "Item %d (%s) tree prefix mismatch", tt.idx, items[tt.idx].Name)
	}
}

func TestBuildTreePrefixMultipleSiblings(t *testing.T) {
	// Test structure: B has two children (C1 and C2)
	//     ┌─ C2 (second child of B, siblingIndex 1)
	//   ┌─┴ C1 (first child of B, siblingIndex 0, has C2 above at depth 3)
	// ┌─┴ B (has C1 above at depth 3)
	// ┌─┴ A (has B above at depth 2)
	// main
	items := []*models.SpiceStackItem{
		{Name: "C2", Depth: 3, SiblingIndex: 1},
		{Name: "C1", Depth: 3, SiblingIndex: 0},
		{Name: "B", Depth: 2, SiblingIndex: 0},
		{Name: "A", Depth: 1, SiblingIndex: 0},
		{Name: "main", Depth: 0, SiblingIndex: 0},
	}

	tests := []struct {
		idx      int
		expected string
	}{
		{0, "    ├──◯ "}, // C2: siblingIndex 1
		{1, "    ┌──◯ "}, // C1: siblingIndex 0
		{2, "  ┌─┴◯ "},   // B
		{3, "┌─┴◯ "},     // A
		{4, ""},          // main: no prefix
	}

	for _, tt := range tests {
		result := buildTreePrefix(items[tt.idx], tt.idx, items)
		result = stripAnsi(result)
		assert.Equal(t, tt.expected, result, "Item %d (%s) tree prefix mismatch", tt.idx, items[tt.idx].Name)
	}
}

func TestBuildTreePrefixNestedSiblings(t *testing.T) {
	// Test structure with siblings at different levels:
	// D is child of C2; C2 and C1 are siblings (children of B); B is child of A
	//     │ ┌─ D (child of C2, ancestor C2 has siblingIndex 1)
	//   ├─┴ C2 (siblingIndex 1, has D above)
	//   ┌─ C1 (siblingIndex 0, no children - D is child of C2)
	// ┌─┴ B (has C1/C2 above)
	// ┌─┴ A (has B above)
	// main
	items := []*models.SpiceStackItem{
		{Name: "D", Depth: 4, SiblingIndex: 0},
		{Name: "C2", Depth: 3, SiblingIndex: 1},
		{Name: "C1", Depth: 3, SiblingIndex: 0},
		{Name: "B", Depth: 2, SiblingIndex: 0},
		{Name: "A", Depth: 1, SiblingIndex: 0},
		{Name: "main", Depth: 0, SiblingIndex: 0},
	}

	tests := []struct {
		idx      int
		expected string
	}{
		{0, "    │ ┌──◯ "}, // D
		{1, "    ├─┴◯ "},   // C2: siblingIndex 1, has D above
		{2, "    ┌──◯ "},   // C1: siblingIndex 0, no children (D is child of C2)
		{3, "  ┌─┴◯ "},     // B
		{4, "┌─┴◯ "},       // A
		{5, ""},            // main: no prefix
	}

	for _, tt := range tests {
		result := buildTreePrefix(items[tt.idx], tt.idx, items)
		result = stripAnsi(result)
		assert.Equal(t, tt.expected, result, "Item %d (%s) tree prefix mismatch", tt.idx, items[tt.idx].Name)
	}
}

func TestBuildTreePrefixSiblingWithCommits(t *testing.T) {
	// Sibling branches where the first has commits should not cause
	// the second sibling to show ┴ (has-children indicator)
	// This tests the fix for commits being incorrectly counted as children
	items := []*models.SpiceStackItem{
		{Name: "sibling1", Depth: 2, SiblingIndex: 0, IsCommit: false},
		{Name: "sibling1", Depth: 3, IsCommit: true, CommitSha: "abc1234", CommitSubject: "Commit on sibling1"},
		{Name: "sibling2", Depth: 2, SiblingIndex: 1, IsCommit: false}, // Should NOT have ┴
		{Name: "parent", Depth: 1, SiblingIndex: 0, IsCommit: false},
		{Name: "main", Depth: 0, SiblingIndex: 0, IsCommit: false},
	}

	result := buildTreePrefix(items[2], 2, items)
	result = stripAnsi(result)

	// sibling2 should have ├──◯ (not ├─┴◯) since it has no children
	// The ┴ should only appear when there are actual child branches, not commits
	assert.NotContains(t, result, "┴", "Sibling branch should not show ┴ when only commits (not branches) are before it")
	assert.Equal(t, "  ├──◯ ", result, "sibling2 should have correct prefix without ┴")
}

// stripAnsi removes ANSI escape codes from a string
func stripAnsi(str string) string {
	// Simple implementation - just remove common ANSI codes
	// For more robust handling, we'd use a regex
	result := strings.Builder{}
	inEscape := false
	for _, r := range str {
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			result.WriteRune(r)
		}
	}
	return result.String()
}
