package gui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Two accessors on one Gui stand for the helpers before and after a repo switch.
func TestRepoListLoadSeqSurvivesHelperRebuild(t *testing.T) {
	g := &Gui{}
	before := &StateAccessor{gui: g}
	after := &StateAccessor{gui: g}

	oldSeq := before.NextRepoListLoadSeq()
	newSeq := after.NextRepoListLoadSeq()

	assert.NotEqual(t, oldSeq, newSeq)
	assert.False(t, before.IsLatestRepoListLoad(oldSeq))
	assert.True(t, before.IsLatestRepoListLoad(newSeq))
}
