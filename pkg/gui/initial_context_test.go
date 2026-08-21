package gui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// normalizeStartupPanel is the validation half of gui.startupPanel (#5877):
// recognized names pass through, everything else (empty, "files" — the
// default the caller already starts on — and typos) normalizes to "" so the
// Files panel stays focused.
func TestNormalizeStartupPanel(t *testing.T) {
	for _, name := range []string{
		"worktrees", "submodules", "branches", "remotes",
		"tags", "commits", "reflog", "stash",
	} {
		assert.Equal(t, name, normalizeStartupPanel(name), name)
	}
	for _, name := range []string{"", "files", "status", "Files", "commit", "nonsense"} {
		assert.Equal(t, "", normalizeStartupPanel(name), name)
	}
}

// The wiring half: initialContext consults the config only when no CLI
// argument selected a panel, and each recognized name selects the matching
// tree field. Source-structural pin (the repo's established pattern for
// context-graph assertions) because building a real ContextTree requires the
// full view harness.
func TestInitialContextWiring(t *testing.T) {
	b, err := os.ReadFile("gui.go")
	assert.NoError(t, err)
	source := string(b)

	assert.Contains(t, source,
		"} else if panel := normalizeStartupPanel(userConfig.Gui.StartupPanel); panel != \"\" {",
		"the config branch must come after the FilterPath/GitArg branches so CLI arguments win")
	assert.Contains(t, source,
		"initialContext(contextTree, startArgs, gui.Config.GetUserConfig())",
		"the caller must pass the user config through")

	for _, pair := range [][2]string{
		{"worktrees", "contextTree.Worktrees"},
		{"submodules", "contextTree.Submodules"},
		{"branches", "contextTree.Branches"},
		{"remotes", "contextTree.Remotes"},
		{"tags", "contextTree.Tags"},
		{"commits", "contextTree.LocalCommits"},
		{"reflog", "contextTree.ReflogCommits"},
		{"stash", "contextTree.Stash"},
	} {
		assert.Contains(t, source,
			"case \""+pair[0]+"\":\n\t\t\tinitialContext = "+pair[1],
			"startupPanel "+pair[0]+" must select "+pair[1])
	}
}
