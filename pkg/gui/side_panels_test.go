package gui

import (
	"sort"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func sortedKeys[V any](m map[string]V) []string {
	keys := lo.Keys(m)
	sort.Strings(keys)
	return keys
}

// The three lookups that translate gui.sidePanels names into views, titles, and
// contexts must each cover exactly the set of valid names, or a config that uses
// a name missing from one of them would hit a nil lookup at runtime.
func TestSidePanelLookupsCoverAllValidTabs(t *testing.T) {
	want := lo.Uniq(config.ValidSidePanelTabs)
	sort.Strings(want)

	gui := NewDummyGui()

	assert.Equal(t, want, sortedKeys(sidePanelViewNames))
	assert.Equal(t, want, sortedKeys(gui.sidePanelTabTitles()))
	assert.Equal(t, want, sortedKeys(sidePanelContexts(gui.contextTree())))
}
